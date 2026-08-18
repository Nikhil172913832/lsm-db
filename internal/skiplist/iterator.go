package skiplist

import "bytes"

type Iterator struct {
	list    *SkipList
	current *Node
}

func (sl *SkipList) NewIterator() *Iterator {
	return &Iterator{
		list: sl,
	}
}

func (it *Iterator) Valid() bool {
	return it.current != nil
}

func (it *Iterator) SeekToFirst() {
	it.current = it.list.head.skipPtrs[0]
}

func (it *Iterator) SeekToLast() {
	temp := it.list.head
	for i := it.list.level - 1; i >= 0; i-- {
		for {
			if temp.skipPtrs[i] != nil {
				temp = temp.skipPtrs[i]
			} else {
				break
			}
		}
	}
	if temp == it.list.head {
		it.current = nil
	} else {
		it.current = temp
	}
}

func (it *Iterator) Seek(target []byte) {
	temp := it.list.head
	for i := it.list.level - 1; i >= 0; i-- {
		for {
			if temp.skipPtrs[i] != nil && bytes.Compare(temp.skipPtrs[i].key, target) < 0 {
				temp = temp.skipPtrs[i]
			} else {
				break
			}
		}
		if temp.skipPtrs[i] != nil && bytes.Equal(temp.skipPtrs[i].key, target) {
			it.current = temp.skipPtrs[i]
			return
		}
	}
	it.current = temp.skipPtrs[0]
}

func (it *Iterator) Next() {
	if it.Valid() {
		it.current = it.current.skipPtrs[0]
	}
}

func (it *Iterator) Key() []byte {
	if it.Valid() {
		return it.current.key
	}
	return nil
}

func (it *Iterator) Value() []byte {
	if it.Valid() {
		return it.current.Value
	}
	return nil
}


