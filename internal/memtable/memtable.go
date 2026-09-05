package memtable

import (
	"sync"

	"github.com/Nikhil172913832/lsm-db/internal/skiplist"
)

type Memtable struct {
	skiplist  *skiplist.SkipList
	sizeBytes int
	threshold int
	mu        sync.RWMutex
}

func NewMemtable(maxLevel int, threshold int) *Memtable {
	return &Memtable{
		skiplist:  skiplist.NewSkipList(maxLevel),
		sizeBytes: 0,
		threshold: threshold,
	}
}

func (mt *Memtable) Put(key, value []byte) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	delta := mt.skiplist.Insert(key, value)
	mt.sizeBytes += delta
}

func (mt *Memtable) Get(key []byte) ([]byte, bool) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	node := mt.skiplist.Search(key)
	if node == nil {
		return nil, false
	}
	return node.Value, true
}

func (mt *Memtable) Size() int {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.sizeBytes
}

func (mt *Memtable) NewIterator() *skiplist.Iterator{
	mt.mu.RLock();
	defer mt.mu.RUnlock()
	return mt.skiplist.NewIterator()
}