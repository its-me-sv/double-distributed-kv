package datastructures

import "github.com/its-me-sv/double-distributed-kv/internal/datastructures/bplustree"

func main() {
	// node = {{"k1": "hi", "k3": "hello"}}
	oldNode := bplustree.Node(make([]byte, bplustree.PAGE_SIZE))
	// leaf node with 2 keys
	oldNode.SetHeader(bplustree.NODE_LEAF, 2)
	// adding k1:hi at 1st pos ie idx 0
	bplustree.NodeAppendKV(oldNode, 0, 0, []byte("k1"), []byte("hi"))
	// adding k2:hello at 2nd pos ie idx 1
	bplustree.NodeAppendKV(oldNode, 1, 0, []byte("k3"), []byte("hello"))

	// add k2 to old node = {{"k1": "hi", "k2": "a", "k3": "hello"}}
	newNode := bplustree.Node(make([]byte, bplustree.PAGE_SIZE))
	newNode.SetHeader(bplustree.NODE_LEAF, 3)
	bplustree.NodeAppendKV(newNode, 0, 0, oldNode.GetKey(0), oldNode.GetVal(0))
	bplustree.NodeAppendKV(newNode, 1, 0, []byte("k2"), []byte("b"))
	bplustree.NodeAppendKV(newNode, 2, 0, oldNode.GetKey(1), oldNode.GetVal(1))

	// remove "k2" = {{"k1": "hi", "k3": "hello"}}
	// new node is now the old node
	oldNode = newNode
	newNode = bplustree.Node(make([]byte, bplustree.PAGE_SIZE))
	newNode.SetHeader(bplustree.NODE_LEAF, 2)
	bplustree.NodeAppendKV(newNode, 0, 0, oldNode.GetKey(0), oldNode.GetVal(0))
	bplustree.NodeAppendKV(newNode, 1, 0, oldNode.GetKey(2), oldNode.GetKey(2))

	// add "a": "b" = {{"k1": "hi", "k3": "hello", "a": "b"}}
	// again new node is now the old node
	oldNode = newNode
	newNode = bplustree.Node(make([]byte, 2*bplustree.PAGE_SIZE))
	newNode.SetHeader(bplustree.NODE_LEAF, 4)
	bplustree.NodeAppendKV(newNode, 0, 0, oldNode.GetKey(0), oldNode.GetVal(0))
	bplustree.NodeAppendKV(newNode, 1, 0, oldNode.GetKey(1), oldNode.GetVal(1))
	bplustree.NodeAppendKV(newNode, 2, 0, []byte("a"), []byte("b"))
}
