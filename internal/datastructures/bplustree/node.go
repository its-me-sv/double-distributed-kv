package bplustree

import (
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
	PAGE_SIZE    = 4096
	MAX_KEY_SIZE = 1000
	MAX_VAL_SIZE = 3000
)

const (
	NODE_INTERNAL = 1
	NODE_LEAF     = 2
)

type Node []byte

// read node type
func (node Node) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// read no. of keys
func (node Node) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// write header
func (node Node) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// read child pointers (internal node)
func (node Node) getPtr(idx uint16) uint16 {
	utils.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint16(node[pos:])
}

// write child pointers (internal node)
func (node Node) setPtr(idx uint16, val uint16) {
	utils.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint16(node[pos:], val)
}

// read offset array (to locate nth key)
func (node Node) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

// write offset value array
func (node Node) setOffset(idx uint16, value uint16) {
	pos := 4*8*node.nkeys() + 2*(idx-1)
	binary.LittleEndian.PutUint16(node[pos:], value)
}

// read position of nth key
func (node Node) kvPos(idx uint16) uint16 {
	utils.Assert(idx <= node.nkeys())
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

// read nth key
func (node Node) GetKey(idx uint16) []byte {
	utils.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:keyLen]
}

// read nth key value (leaf nodes)
func (node Node) GetVal(idx uint16) []byte {
	utils.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	valueLen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+keyLen:][:valueLen]
}

// node size (in bytes)
func (node Node) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

// idx - position of key or value or pointer
// ptr - nth child pointer (unused for leaf nodes)
// key, val - the KV pair (empty for internal nodes)
// assumption - keys are added in order because uses previous
// values of the offsets array to locate the KV
func NodeAppendKV(node Node, idx uint16, ptr uint16, key []byte, val []byte) {
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
