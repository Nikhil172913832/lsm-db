package engine

import (
	"fmt"
	"path/filepath"
)

func (e *Engine) maybeFlush() error {
	if e.activeMemtable.SizeBytes > e.memtableThreshold {
		e.mu.Lock()
		e.immutableWAL = e.activeWAL
		e.immutableMemtable = e.activeMemtable
		filename := fmt.Sprintf("%s_%06d", "wal", e.nextSeqNum)
		filePath := filepath.Join(e.walDir, filename)
		var err error
		e.activeWAL, e.activeMemtable, err = LoadWALAndMemtable(filePath, e.maxLevel, e.memtableThreshold)
		if err != nil {
			e.activeWAL = e.immutableWAL
			e.activeMemtable = e.immutableMemtable
			e.immutableWAL = nil
			e.immutableMemtable = nil
			e.mu.Unlock()
			return err
		}
		e.mu.Unlock()
	}

	return nil
}
