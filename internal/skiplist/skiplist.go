package skiplist

import (
	"bytes"
)

type SkipList struct {
	head     *Node
	maxLevel int
	level    int
}

func New(maxLevel int) *SkipList {
	head := &Node{
		key:      nil,
		Value:    nil,
		skipPtrs: make([]*Node, maxLevel),
	}
	return &SkipList{
		head:     head,
		maxLevel: maxLevel,
		level:    0,
	}
}

func (sl *SkipList) Insert(key, value []byte) int {

	temp := sl.head
	level := sl.level - 1
	var prevPtr []*Node
	for {
		if temp == nil {
			break
		}
		change := false
		for ; level >= 0; level-- {
			if temp.skipPtrs[level] != nil {
				comp := bytes.Compare(temp.skipPtrs[level].key, key)
				if comp < 0 {
					temp = temp.skipPtrs[level]
					change = true
					break
				} else if comp == 0 {
					delta := len(value) - len(temp.skipPtrs[level].Value)
					temp.skipPtrs[level].Value = value
					return delta
				}
			}
			prevPtr = append(prevPtr, temp)
		}
		if !change {
			break
		}
	}
	newLevel := RandomLevel(sl.maxLevel)
	newNode := Node{key: key, Value: value}
	for i := len(prevPtr) - 1; i >= max(len(prevPtr)-newLevel, 0); i-- {
		newNode.skipPtrs = append(newNode.skipPtrs, prevPtr[i].skipPtrs[len(prevPtr)-1-i])
		prevPtr[i].skipPtrs[len(prevPtr)-1-i] = &newNode
	}
	for i := len(prevPtr); i < newLevel; i++ {
		sl.head.skipPtrs[i] = &newNode
		newNode.skipPtrs = append(newNode.skipPtrs, nil)
	}
	sl.level = max(sl.level, newLevel)
	return len(value) + len(key)
}

func (sl *SkipList) Search(key []byte) *Node {
	temp := sl.head
	level := sl.level - 1
	for {
		if temp == nil {
			break
		}
		change := false
		for ; level >= 0; level-- {
			if temp.skipPtrs[level] != nil {
				comp := bytes.Compare(temp.skipPtrs[level].key, key)
				if comp < 0 {
					temp = temp.skipPtrs[level]
					change = true
					break
				} else if comp == 0 {
					return temp.skipPtrs[level]
				}
			}
		}
		if !change {
			break
		}
	}
	return nil
}
