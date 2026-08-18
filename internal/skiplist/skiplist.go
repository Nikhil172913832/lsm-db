package skiplist

import (
	"bytes"
)

type SkipList struct {
	head     *Node
	maxLevel int
	level    int
}

func NewSkipList(maxLevel int) *SkipList {
	head := &Node{
		skipPtrs: make([]*Node, maxLevel),
	}
	return &SkipList{
		head:     head,
		maxLevel: maxLevel,
	}
}

func (sl *SkipList) Insert(key, value []byte) int {
	temp := sl.head
	prevPtr := make([]*Node, sl.maxLevel)
	for i := sl.level - 1; i >= 0; i-- {
		for {
			if temp.skipPtrs[i] != nil && bytes.Compare(temp.skipPtrs[i].key, key) < 0 {
				temp = temp.skipPtrs[i]
			} else {
				break
			}
		}
		if temp.skipPtrs[i] != nil && bytes.Equal(temp.skipPtrs[i].key, key) {
			delta := len(value) - len(temp.skipPtrs[i].Value)
			value = bytes.Clone(value)
			temp.skipPtrs[i].Value = value
			return delta
		}
		prevPtr[i] = temp
	}
	newLevel := RandomLevel(sl.maxLevel)
	key = bytes.Clone(key)
    value = bytes.Clone(value)
	newNode := Node{key: key, Value: value}
	newNode.skipPtrs = make([]*Node, newLevel)
	for i := sl.level; i < newLevel; i++ {
		prevPtr[i] = sl.head
	}
	for i := range newLevel {
		newNode.skipPtrs[i] = prevPtr[i].skipPtrs[i]
		prevPtr[i].skipPtrs[i] = &newNode
	}
	sl.level = max(sl.level, newLevel)
	return len(value) + len(key)
}

func (sl *SkipList) Search(key []byte) *Node {
	temp := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for {
			if temp.skipPtrs[i] != nil && bytes.Compare(temp.skipPtrs[i].key, key) < 0 {
				temp = temp.skipPtrs[i]
			} else {
				break
			}
		}
		if temp.skipPtrs[i] != nil && bytes.Equal(temp.skipPtrs[i].key, key) {
			return temp.skipPtrs[i]
		}
	}
	return nil
}
