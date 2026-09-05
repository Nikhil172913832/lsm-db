package engine

import (
	"fmt"
	"path/filepath"
)

func (e *Engine) maybeFlush() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.state.current.Size() > e.memtableThreshold {

		if e.state.current.Size() <= e.memtableThreshold {
			return nil
		}
		filename := fmt.Sprintf("%s_%06d", "wal", e.nextSeqNum)
		filePath := filepath.Join(e.walDir, filename)
		newCurrent, err := NewWriteState(filePath, e.maxLevel, e.memtableThreshold, e.maxKeySize, e.maxValueSize)
		if err != nil {
			return err
		}
		newImmutable := append(e.state.immutable, e.state.current)
		e.state.current.Freeze()
		e.state = &EngineState{
			current:   newCurrent,
			immutable: newImmutable,
			// ssTables: e.ssTables TO DO
		}
		e.nextSeqNum++
	}
	return nil
}
