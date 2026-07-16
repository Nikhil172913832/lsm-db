package memtable

import (
	"github.com/Nikhil172913832/lsm-db/internal/skiplist"
)

type Memtable struct {
	skiplist  *skiplist.SkipList
	sizeBytes int
	threshold int
}

func New(maxLevel int, threshold int) *Memtable {
	return &Memtable{
		skiplist:  skiplist.New(maxLevel),
		sizeBytes: 0,
		threshold: threshold,
	}
}

func (mt *Memtable) Put(key, value []byte) error {
	if len(value) == 0 {
		return ErrEmptyOrNilValue
	}
	delta := mt.skiplist.Insert(key, value)
	mt.sizeBytes += delta
	return nil
}

func (mt *Memtable) Get(key []byte) ([]byte, bool) {
	node := mt.skiplist.Search(key)
	if node == nil {
		return nil, false
	}
	return node.Value, true
}

func (mt *Memtable) Delete(key []byte) {
	delta := mt.skiplist.Insert(key, nil)
	mt.sizeBytes += delta
}

func (mt *Memtable) Size() int{
	return mt.sizeBytes
}
