package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
}

func NewEngine(walDir, sstableDir string, maxLevel, memtableThreshold int) (*Engine, error) {
	engine := Engine{
		walDir:            walDir,
		ssTableDir:        sstableDir,
		maxLevel:          maxLevel,
		memtableThreshold: memtableThreshold,
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
		state.current, err = NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold)
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
			immutable, err := NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold)
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
		state.current, err = NewWriteState(filePath, engine.maxLevel, engine.memtableThreshold)
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
