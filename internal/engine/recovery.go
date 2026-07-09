package engine

import (
	"github.com/Nikhil172913832/lsm-db/internal/memtable"
	"github.com/Nikhil172913832/lsm-db/internal/wal"
)

func LoadWALAndMemtable(path string, maxLevel, threshold int) (*wal.WAL, *memtable.Memtable, error) {
	wal, err := wal.New(path)
	if err != nil {
		return nil, nil, err
	}
	memtable := memtable.New(maxLevel, threshold)
	if err := wal.ReadAll(ReplayInto(memtable)); err != nil {
		return nil, nil, err
	}
	return wal, memtable, nil
}

func ReplayInto(memtable *memtable.Memtable) func(key, value []byte) error {
	return func(key, value []byte) error {
		if len(value) == 0 {
			memtable.Delete(key)
			return nil
		}
		return memtable.Put(key, value)
	}
}
