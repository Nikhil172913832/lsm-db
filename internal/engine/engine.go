package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Engine struct {
	mu        sync.RWMutex
	active    *State
	immutable *State
	// sstables          []*SSTable #TO DO
	walDir            string
	ssTableDir        string
	maxLevel          int
	memtableThreshold int
	nextSeqNum        int
}

func New(walDir, sstableDir string, maxLevel, memtableThreshold int) (*Engine, error) {
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
	if len(entries) == 0 {
		engine.nextSeqNum = 1
		filename := fmt.Sprintf("%s_%06d", "wal", engine.nextSeqNum)
		filePath := filepath.Join(walDir, filename)
		engine.active.LoadWALAndMemtable(filePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
	} else if len(entries) == 1 {
		filePath := filepath.Join(walDir, entries[0].Name())
		fileNameParts := strings.Split(entries[0].Name(), "_")
		if len(fileNameParts) != 2 {
			return nil, InvalidWALFileName
		}
		engine.nextSeqNum, err = strconv.Atoi(fileNameParts[1])
		if err != nil {
			return nil, err
		}
		engine.nextSeqNum++
		engine.active.LoadWALAndMemtable(filePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
	} else if len(entries) == 2 {
		activeFilePath := filepath.Join(walDir, entries[len(entries)-1].Name())
		fileNameParts := strings.Split(entries[len(entries)-1].Name(), "_")
		if len(fileNameParts) != 2 {
			return nil, InvalidWALFileName
		}
		engine.nextSeqNum, err = strconv.Atoi(fileNameParts[1])
		if err != nil {
			return nil, err
		}
		engine.nextSeqNum++
		engine.active.LoadWALAndMemtable(activeFilePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
		immutableFilePath := filepath.Join(walDir, entries[0].Name())
		engine.immutable.LoadWALAndMemtable(immutableFilePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, InvalidNoOfFilesInWALDir
	}
	return &engine, nil
}
