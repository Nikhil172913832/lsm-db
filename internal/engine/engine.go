package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Nikhil172913832/lsm-db/internal/base"
)

type EngineState struct {
	current   *WriteState
	immutable []*WriteState
	// ssTables          []*SSTable #TO DO
}
type Engine struct {
	mu                sync.RWMutex
	state             *EngineState
	walDir            string
	ssTableDir        string
	maxLevel          int
	memtableThreshold int
	nextSeqNum        int
	maxKeySize        uint32
	maxValueSize      uint32
}

func NewEngine(walDir, sstableDir string, maxLevel, memtableThreshold int, maxKeySize, maxValueSize uint32) (*Engine, error) {
	engine := Engine{
		walDir:            walDir,
		ssTableDir:        sstableDir,
		maxLevel:          maxLevel,
		memtableThreshold: memtableThreshold,
		maxKeySize:        maxKeySize,
		maxValueSize:      maxValueSize,
	}
	if _, err := os.Stat(engine.walDir); err != nil {
		return nil, err
	}
	if _, err := os.Stat(engine.ssTableDir); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(engine.walDir)
	if err != nil {
		return nil, err
	}
	state := &EngineState{}
	if len(entries) == 0 {
		engine.nextSeqNum = 1
		filename := fmt.Sprintf("%s_%06d", "wal", engine.nextSeqNum)
		filePath := filepath.Join(walDir, filename)
		state.current, err = NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold, engine.maxKeySize, engine.maxValueSize)
		if err != nil {
			return nil, err
		}
	} else {
		for i := 0; i < len(entries)-1; i++ {
			filePath := filepath.Join(walDir, entries[i].Name())
			fileNameParts := strings.Split(entries[i].Name(), "_")
			if len(fileNameParts) != 2 {
				return nil, InvalidWALFileName
			}
			immutable, err := NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold, engine.maxKeySize, engine.maxValueSize)
			if err != nil {
				return nil, err
			}
			immutable.Freeze()
			state.immutable = append(state.immutable, immutable)
		}
		filePath := filepath.Join(walDir, entries[len(entries)-1].Name())
		fileNameParts := strings.Split(entries[len(entries)-1].Name(), "_")
		if len(fileNameParts) != 2 {
			return nil, InvalidWALFileName
		}
		state.current, err = NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold, engine.maxKeySize, engine.maxValueSize)
		if err != nil {
			return nil, err
		}
		engine.nextSeqNum, err = strconv.Atoi(fileNameParts[1])
		if err != nil {
			return nil, err
		}
	}
	engine.nextSeqNum++
	engine.state = state
	return &engine, nil
}

func (e *Engine) write(OpType byte, key, value []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if uint32(len(key)) > e.maxKeySize {
		return ErrKeySizeExceedMaxLimit
	}
	if uint32(len(value)) > e.maxValueSize {
		return ErrValSizeExceedMaxLimit
	}
	var st *WriteState
	for {
		e.mu.RLock()
		st = e.state.current
		e.mu.RUnlock()
		if st.AcquireForWrite() {
			break
		}

	}
	err := st.wal.Append(OpType, key, value)
	if err != nil {
		st.ReleaseWrite()
		return err
	}
	st.mem.Put(key, value)
	st.ReleaseWrite()
	e.maybeFlush()
	return nil
}

func (e *Engine) Put(key, value []byte) error {
	if value == nil {
		value = []byte{}
	}
	return e.write(base.OpPut, key, value)
}

func (e *Engine) Delete(key []byte) error {
	return e.write(base.OpDelete, key, nil)
}

func (e *Engine) Get(key []byte) ([]byte, bool, error) {
	if len(key) == 0 {
		return nil, false, ErrEmptyKey
	}
	e.mu.RLock()
	currentSt := e.state.current
	immutableSts := e.state.immutable
	e.mu.RUnlock()
	val, found := currentSt.mem.Get(key)
	if found && val != nil {
		return val, true, nil
	}
	if found && val == nil {
		return val, false, nil
	}
	for i := len(immutableSts) - 1; i >= 0; i-- {
		val, found := immutableSts[i].mem.Get(key)
		if found && val != nil {
			return val, true, nil
		}
		if found && val == nil {
			return val, false, nil
		}
	}
	return nil, false, nil
}

func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	err := e.state.current.wal.Close()
	if err != nil{
		return err
	}
	for _, state := range e.state.immutable{
		err = state.wal.Close()
		if(err != nil && !errors.Is(err, os.ErrClosed)){
			return err
		}
	}
	return nil
}
