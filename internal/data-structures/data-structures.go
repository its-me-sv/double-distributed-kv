package datastructures

func main() {
	// node = {{"k1": "hi", "k3": "hello"}}
	oldNode := BNode(make([]byte, BTREE_PAGE_SIZE))
	// leaf node with 2 keys
	oldNode.setHeader(BNODE_LEAF, 2)
	// adding k1:hi at 1st pos ie idx 0
	bNodeAppendKV(oldNode, 0, 0, []byte("k1"), []byte("hi"))
	// adding k2:hello at 2nd pos ie idx 1
	bNodeAppendKV(oldNode, 1, 0, []byte("k3"), []byte("hello"))

	// add k2 to old node = {{"k1": "hi", "k2": "a", "k3": "hello"}}
	newNode := BNode(make([]byte, BTREE_PAGE_SIZE))
	newNode.setHeader(BNODE_LEAF, 3)
	bNodeAppendKV(newNode, 0, 0, oldNode.getKey(0), oldNode.getVal(0))
	bNodeAppendKV(newNode, 1, 0, []byte("k2"), []byte("b"))
	bNodeAppendKV(newNode, 2, 0, oldNode.getKey(1), oldNode.getVal(1))

	// remove "k2" = {{"k1": "hi", "k3": "hello"}}
	// new node is now the old node
	oldNode = newNode
	newNode = BNode(make([]byte, BTREE_PAGE_SIZE))
	newNode.setHeader(BNODE_LEAF, 2)
	bNodeAppendKV(newNode, 0, 0, oldNode.getKey(0), oldNode.getVal(0))
	bNodeAppendKV(newNode, 1, 0, oldNode.getKey(2), oldNode.getKey(2))

	// add "a": "b" = {{"k1": "hi", "k3": "hello", "a": "b"}}
	// again new node is now the old node
	oldNode = newNode
	newNode = BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	newNode.setHeader(BNODE_LEAF, 4)
	bNodeAppendKV(newNode, 0, 0, oldNode.getKey(0), oldNode.getVal(0))
	bNodeAppendKV(newNode, 1, 0, oldNode.getKey(1), oldNode.getVal(1))
	bNodeAppendKV(newNode, 2, 0, []byte("a"), []byte("b"))
}
