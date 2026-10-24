package datastructures

import (
	"bytes"
	"encoding/binary"

	"github.com/its-me-sv/double-distributed-kv/internal/utils"
)

// node format
// | type | nkeys |  pointers  |  offsets   | key-values | unused |
// |  2B  |   2B  | nkeys × 8B | nkeys × 2B |     ...    |        |

// For example, a leaf node {"k1":"hi", "k3":"hello"} is encoded as:

// | type | nkeys | pointers | offsets |            key-values           | unused |
// |   2  |   2   | nil nil  |  8 19   | 2 2 "k1" "hi"  2 5 "k3" "hello" |        |
// |  2B  |  2B   |   2×8B   |  2×2B   | 4B + 2B + 2B + 4B + 2B + 5B     |        |

const (
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

const (
	BNODE_INTERNAL = 1
	BNODE_LEAF     = 2
)

type BNode []byte

// read node type
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// read no. of keys
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// write header
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// read child pointers (internal node)
func (node BNode) getPtr(idx uint16) uint16 {
	utils.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint16(node[pos:])
}

// write child pointers (internal node)
func (node BNode) setPtr(idx uint16, val uint16) {
	utils.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint16(node[pos:], val)
}

// read offset array (to locate nth key)
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

// write offset value array
func (node BNode) setOffset(idx uint16, value uint16) {
	pos := 4*8*node.nkeys() + 2*(idx-1)
	binary.LittleEndian.PutUint16(node[pos:], value)
}

// read position of nth key
func (node BNode) kvPos(idx uint16) uint16 {
	utils.Assert(idx <= node.nkeys())
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

// read nth key
func (node BNode) getKey(idx uint16) []byte {
	utils.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:keyLen]
}

// read nth key value (leaf nodes)
func (node BNode) getVal(idx uint16) []byte {
	utils.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	valueLen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+keyLen:][:valueLen]
}

// node size (in bytes)
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

// idx - position of key or value or pointer
// ptr - nth child pointer (unused for leaf nodes)
// key, val - the KV pair (empty for internal nodes)
// assumption - keys are added in order because uses previous
// values of the offsets array to locate the KV
func bNodeAppendKV(node BNode, idx uint16, ptr uint16, key []byte, val []byte) {
	// pointers
	node.setPtr(idx, ptr)

	// KVs
	pos := node.kvPos(idx)

	// data sizes
	keyLen := uint16(len(key))
	valLen := uint16(len(val))
	binary.LittleEndian.PutUint16(node[pos+0:], keyLen)
	binary.LittleEndian.PutUint16(node[pos+2:], valLen)

	// writing data
	copy(node[pos+4:], key)
	copy(node[pos+4+keyLen:], val)

	// offset updation for next key
	node.setOffset(idx+1, node.getOffset(idx)+4+keyLen+valLen)
}

// inserting new key at position idx
func bLeafInsert(new, old BNode, idx uint16, key, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	// copy keys before index
	bNodeAppendRange(new, old, 0, 0, idx)
	// insert the key
	bNodeAppendKV(new, idx, 0, key, val)
	// copy keys after index
	bNodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// updating existing key's value if exists
func bLeafUpdate(new, old BNode, idx uint16, key, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys())
	bNodeAppendRange(new, old, 0, 0, idx)
	bNodeAppendKV(new, idx, 0, key, val)
	bNodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-(idx+1))
}

// insert or update based on key availability using lookup
func bLeafSet(new, old BNode, key, val []byte) {
	idx := bNodeLookupLS(old, key)
	if bytes.Equal(key, old.getKey(idx)) {
		bLeafUpdate(new, old, idx, key, val)
	} else {
		bLeafInsert(new, old, idx+1, key, val)
	}
}

// use to find key or closet pos to key to maintain sortorder
func bNodeLookupLS(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()

	var i uint16
	for i = 0; i < nkeys; i += 1 {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp == 0 {
			return i
		}
		if cmp > 0 {
			return i - 1
		}
	}

	return i - 1
}

// copying multiple KVs and pointers into position
func bNodeAppendRange(new, old BNode, dstNew, srcOld uint16, n uint16) {
	for i := uint16(0); i < n; i += 1 {
		dst, src := dstNew+i, srcOld+i
		bNodeAppendKV(new, dst, old.getPtr(src), old.getKey(src), old.getVal(src))
	}
}

// splitting oversized nodes into 2
func bNodeSplitInto2(left, right, old BNode) {
	utils.Assert(old.nkeys() >= 2)
	// initial size
	nLeft := old.nkeys() / 2

	// fitting left half
	leftBytes := func() uint16 {
		return 4 + 8*nLeft + 2*nLeft + old.getOffset(nLeft)
	}
	for leftBytes() > BTREE_PAGE_SIZE {
		nLeft -= 1
	}
	utils.Assert(nLeft >= 1)

	// fitting right half
	rightBytes := func() uint16 {
		return old.nbytes() - leftBytes() + 4
	}
	for rightBytes() > BTREE_PAGE_SIZE {
		nLeft += 1
	}
	utils.Assert(nLeft < old.nkeys())
	nRight := old.nkeys() - nLeft

	// new nodes
	left.setHeader(old.btype(), nLeft)
	right.setHeader(old.btype(), nRight)
	bNodeAppendRange(left, old, 0, 0, nLeft)
	bNodeAppendRange(right, old, 0, nLeft, nRight)

	// left half could still be too big
	utils.Assert(right.nbytes() <= BTREE_PAGE_SIZE)
}

// splitting too large oversized nodes into 3
func bNodeSplitInto3(old BNode) (uint16, [3]BNode) {
	// no splitting - 1 node
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}

	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	right := BNode(make([]byte, BTREE_PAGE_SIZE))
	bNodeSplitInto2(left, right, old)

	// not too large but just oversized - 2 nodes
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left := left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}

	// when too large oversized - 3 nodes
	leftOfLeft := BNode(make([]byte, BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, BTREE_PAGE_SIZE))
	bNodeSplitInto2(leftOfLeft, middle, left)
	utils.Assert(leftOfLeft.nbytes() <= BTREE_PAGE_SIZE)

	return 3, [3]BNode{leftOfLeft, middle, right}
}
