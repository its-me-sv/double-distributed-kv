package bplustree

type Tree struct {
	// non zero page number
	root uint64
	// read data by page number
	get func(uint64) []byte
	// allocate new page
	new func([]byte) uint64
	// deallocate by page number
	del func(uint64)
}
