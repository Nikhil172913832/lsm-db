package engine

import (
	"github.com/Nikhil172913832/lsm-db/internal/memtable"
	"github.com/Nikhil172913832/lsm-db/internal/wal"
)

type State struct{
	wal *wal.WAL
	mem *memtable.Memtable
}

func (st *State) LoadWALAndMemtable(path string, maxLevel, threshold int) error {
	wal, err := wal.New(path)
	if err != nil {
		return err
	}
	mem := memtable.New(maxLevel, threshold)
	if err := wal.ReadAll(ReplayInto(mem)); err != nil {
		return err
	}
	st.wal = wal
	st.mem = mem
	return nil
}