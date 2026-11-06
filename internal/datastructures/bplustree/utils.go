package bplustree

import (
	"bytes"

	"github.com/its-me-sv/double-distributed-kv/internal/utils"
)

// inserting new key at position idx
func LeafInsert(new, old Node, idx uint16, key, val []byte) {
	new.SetHeader(NODE_LEAF, old.nkeys()+1)
	// copy keys before index
	NodeAppendRange(new, old, 0, 0, idx)
	// insert the key
	NodeAppendKV(new, idx, 0, key, val)
	// copy keys after index
	NodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// updating existing key's value if exists
func LeafUpdate(new, old Node, idx uint16, key, val []byte) {
	new.SetHeader(NODE_LEAF, old.nkeys())
	NodeAppendRange(new, old, 0, 0, idx)
	NodeAppendKV(new, idx, 0, key, val)
	NodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-(idx+1))
}

// insert or update based on key availability using lookup
func LeafSet(new, old Node, key, val []byte) {
	idx := NodeLookupLS(old, key)
	if bytes.Equal(key, old.GetKey(idx)) {
		LeafUpdate(new, old, idx, key, val)
	} else {
		LeafInsert(new, old, idx+1, key, val)
	}
}

// use to find key or closet pos to key to maintain sortorder
func NodeLookupLS(node Node, key []byte) uint16 {
	nkeys := node.nkeys()

	var i uint16
	for i = 0; i < nkeys; i += 1 {
		cmp := bytes.Compare(node.GetKey(i), key)
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
func NodeAppendRange(new, old Node, dstNew, srcOld uint16, n uint16) {
	for i := uint16(0); i < n; i += 1 {
		dst, src := dstNew+i, srcOld+i
		NodeAppendKV(new, dst, old.getPtr(src), old.GetKey(src), old.GetVal(src))
	}
}

// splitting oversized nodes into 2
func NodeSplitInto2(left, right, old Node) {
	utils.Assert(old.nkeys() >= 2)
	// initial size
	nLeft := old.nkeys() / 2

	// fitting left half
	leftBytes := func() uint16 {
		return 4 + 8*nLeft + 2*nLeft + old.getOffset(nLeft)
	}
	for leftBytes() > PAGE_SIZE {
		nLeft -= 1
	}
	utils.Assert(nLeft >= 1)

	// fitting right half
	rightBytes := func() uint16 {
		return old.nbytes() - leftBytes() + 4
	}
	for rightBytes() > PAGE_SIZE {
		nLeft += 1
	}
	utils.Assert(nLeft < old.nkeys())
	nRight := old.nkeys() - nLeft

	// new nodes
	left.SetHeader(old.btype(), nLeft)
	right.SetHeader(old.btype(), nRight)
	NodeAppendRange(left, old, 0, 0, nLeft)
	NodeAppendRange(right, old, 0, nLeft, nRight)

	// left half could still be too big
	utils.Assert(right.nbytes() <= PAGE_SIZE)
}

// splitting too large oversized nodes into 3
func NodeSplitInto3(old Node) (uint16, [3]Node) {
	// no splitting - 1 node
	if old.nbytes() <= PAGE_SIZE {
		old = old[:PAGE_SIZE]
		return 1, [3]Node{old}
	}

	left := Node(make([]byte, 2*PAGE_SIZE))
	right := Node(make([]byte, PAGE_SIZE))
	NodeSplitInto2(left, right, old)

	// not too large but just oversized - 2 nodes
	if left.nbytes() <= PAGE_SIZE {
		left := left[:PAGE_SIZE]
		return 2, [3]Node{left, right}
	}

	// when too large oversized - 3 nodes
	leftOfLeft := Node(make([]byte, PAGE_SIZE))
	middle := Node(make([]byte, PAGE_SIZE))
	NodeSplitInto2(leftOfLeft, middle, left)
	utils.Assert(leftOfLeft.nbytes() <= PAGE_SIZE)

	return 3, [3]Node{leftOfLeft, middle, right}
}
