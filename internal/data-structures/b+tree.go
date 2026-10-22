package datastructures

import (
	"encoding/binary"
)

type Node struct {
	keys [][]byte
	vals [][]byte // for leaf node
	kids []*Node  // for internal node
}

type BNode []byte

// node type
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// no. of keys
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// read child pointers
func (node BNode) getPtr(idx uint16) uint16 {
	assert(idx < node.nkeys())
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint16(node[pos:])
}

// write child pointers
func (node BNode) setPtr(idx uint16, val uint16) {
	assert(idx < node.nkeys())
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint16(node[pos:], val)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:keyLen]
}

func (node BNode) getVal(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	keyLen := binary.LittleEndian.Uint16(node[pos:])
	valueLen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+keyLen:][:valueLen]
}

func Encode(node *Node) []byte

func Decode(page []byte) (*Node, error)

type NodeType int

const (
	NODE_INTERNAL NodeType = iota + 1
	NODE_LEAF
)

const (
	TREE_PAGE_SIZE    = 4096
	TREE_MAX_KEY_SIZE = 1000
	TREE_MAX_VAL_SIZE = 3000
)

func init() {
	node1max := 4 + 1*8 + 1*2 + 4 + TREE_MAX_KEY_SIZE + TREE_MAX_VAL_SIZE
	assert(node1max <= TREE_PAGE_SIZE)
}

func assert(b bool) {
	if !b {
		panic("NOT OK")
	}
}
