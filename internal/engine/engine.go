package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Nikhil172913832/lsm-db/internal/memtable"
	"github.com/Nikhil172913832/lsm-db/internal/wal"
)

type Engine struct {
	mu                sync.RWMutex
	activeMemtable    *memtable.Memtable
	immutableMemtable *memtable.Memtable
	activeWAL         *wal.WAL
	immutableWAL      *wal.WAL
	// sstables          []*SSTable
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
		engine.activeWAL, engine.activeMemtable, err = LoadWALAndMemtable(filePath, engine.maxLevel, engine.memtableThreshold)
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
		engine.activeWAL, engine.activeMemtable, err = LoadWALAndMemtable(filePath, engine.maxLevel, engine.memtableThreshold)
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
		engine.activeWAL, engine.activeMemtable, err = LoadWALAndMemtable(activeFilePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
		immutableFilePath := filepath.Join(walDir, entries[0].Name())
		engine.immutableWAL, engine.immutableMemtable, err = LoadWALAndMemtable(immutableFilePath, engine.maxLevel, engine.memtableThreshold)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, InvalidNoOfFilesInWALDir
	}
	return &engine, nil
}
