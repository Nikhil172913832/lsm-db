package engine

import (
	"fmt"
	"path/filepath"
)

func (e *Engine) maybeFlush() error {
		if e.state.current.Size() > e.memtableThreshold {
			e.mu.Lock()
			defer e.mu.Unlock()
			filename := fmt.Sprintf("%s_%06d", "wal", e.nextSeqNum)
			filePath := filepath.Join(e.walDir, filename)
			newCurrent, err := NewWriteState(filePath, e.maxLevel, e.memtableThreshold)
			newImmutable := append(e.state.immutable, e.state.current)
			if err != nil{
				return err
			}
			e.state = &EngineState{
				current: newCurrent,
				immutable: newImmutable,
				// ssTables: e.ssTables TO DO
			}
			e.nextSeqNum++
		}
		return nil
}
