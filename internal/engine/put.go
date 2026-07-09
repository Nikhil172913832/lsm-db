package engine

func (e *Engine) Put(key, value []byte) error {
	if err := e.activeWAL.Append(key, value); err != nil {
		return err
	}
	if err := e.activeMemtable.Put(key, value); err != nil {
		return err
	}
	if err := e.maybeFlush(); err != nil {
		return err
	}
	return nil
}
