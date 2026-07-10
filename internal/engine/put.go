package engine

func (e *Engine) Put(key, value []byte) error {
	if err := e.active.wal.Append(key, value); err != nil {
		return err
	}
	if err := e.active.mem.Put(key, value); err != nil {
		return err
	}
	if err := e.maybeFlush(); err != nil {
		return err
	}
	return nil
}
