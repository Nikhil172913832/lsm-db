package wal

import "encoding/binary"

func (w *WAL) readRecord(offset int64) ([]byte, []byte, int64, error) {
	start := offset
	cursor := offset
	var lenBuf [4]byte
	if _, err := w.file.ReadAt(lenBuf[:], cursor); err != nil {
		return nil, nil, start, err
	}
	keyLen := binary.LittleEndian.Uint32(lenBuf[:])
	if keyLen > MaxKeySize {
		return nil, nil, start, ErrCorruptRecord
	}
	cursor += 4
	key := make([]byte, keyLen)
	if _, err := w.file.ReadAt(key, cursor); err != nil {
		return nil, nil, start, err
	}
	cursor += int64(keyLen)
	if _, err := w.file.ReadAt(lenBuf[:], cursor); err != nil {
		return nil, nil, start, err
	}
	valLen := binary.LittleEndian.Uint32(lenBuf[:])
	if valLen > MaxValueSize {
		return nil, nil, start, ErrCorruptRecord
	}
	cursor += 4
	var val []byte
	if valLen > 0 {
		val = make([]byte, valLen)
		if _, err := w.file.ReadAt(val, cursor); err != nil {
			return nil, nil, start, err
		}
		cursor += int64(valLen)
	}
	return key, val, cursor, nil
}
