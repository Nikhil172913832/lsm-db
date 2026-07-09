package skiplist

type Node struct {
	key      []byte
	Value    []byte
	skipPtrs []*Node
}
