package engine

import (
	"sync"

	"github.com/Nikhil172913832/lsm-db/internal/memtable"
	"github.com/Nikhil172913832/lsm-db/internal/wal"
)

type WriteState struct {
	wal     *wal.WAL
	mem     *memtable.Memtable
	mu      sync.Mutex
	frozen  bool
	writers int
}

func NewWriteState(path string, maxLevel, threshold int) (*WriteState, error) {
	wal, err := wal.New(path)
	if err != nil {
		return nil, err
	}
	mem := memtable.NewMemtable(maxLevel, threshold)
	if err := wal.ReadAll(ReplayInto(mem)); err != nil {
		return nil, err
	}
	return &WriteState{
		wal: wal,
		mem: mem,
	}, nil
}

func ReplayInto(memtable *memtable.Memtable) func(key, value []byte) error {
	return func(key, value []byte) error {
		if len(value) == 0 {
			memtable.Delete(key)
			return nil
		}
		memtable.Put(key, value)
		return nil
	}
}

func (st *WriteState) AcquireForWrite() bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.frozen {
		return false
	}
	st.writers++
	return true
}

func (st *WriteState) ReleaseWrite() (becameIdle bool, err error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.writers == 0 {
		return false, ErrActiveStateDoubleFree
	}
	st.writers--
	if st.writers == 0 {
		return true, nil
	}
	return false, nil
}

func (st *WriteState) Freeze() bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.frozen {
		return false
	}
	st.frozen = true
	return true
}

func (st *WriteState) WalPath() string {
	return st.wal.Path()
}

func (st *WriteState) Size() int {
	return st.mem.Size()
}
