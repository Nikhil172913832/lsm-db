package wal

import (
	"errors"
	"io"
)

func (w *WAL) ReadAll(fn func(key, value []byte) error) error {
	var offset int64 = 0
	for {
		var key, val []byte
		var err error
		key, val, offset, err = w.readRecord(offset)
		switch {
		case err == nil:
			if err := fn(key, val); err != nil {
				return err
			}
		case errors.Is(err, io.EOF):
			return nil
		case errors.Is(err, io.ErrUnexpectedEOF):
			if err := w.file.Truncate(offset); err != nil {
				return err
			}
			return nil
		default:
			return err
		}
	}
}
