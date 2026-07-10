package engine

import (
	"fmt"
	"path/filepath"
)

func (e *Engine) maybeFlush() error {
	if e.active.mem.SizeBytes > e.memtableThreshold {
		e.mu.Lock()
		e.immutable = e.active
		filename := fmt.Sprintf("%s_%06d", "wal", e.nextSeqNum)
		filePath := filepath.Join(e.walDir, filename)
		var err error
		e.active.LoadWALAndMemtable(filePath, e.maxLevel, e.memtableThreshold)
		if err != nil {
			e.active = e.immutable
			e.immutable = nil
			e.mu.Unlock()
			return err
		}
		e.mu.Unlock()
	}

	return nil
}
