package wal

import (
	"bytes"
	"encoding/binary"
	"os"
	"sync"
)

type WAL struct {
	file *os.File
	path string
	mu   sync.Mutex
}

func New(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{
		file: file,
		path: path,
	}, nil
}

func (w *WAL) Append(key, val []byte) error {
	if uint32(len(key)) > MaxKeySize {
		return ErrKeySizeExceedMaxLimit
	}
	if uint32(len(val)) > MaxValueSize {
		return ErrValSizeExceedMaxLimit
	}
	var buf bytes.Buffer
	var lenBuf [4]byte
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(len(key)))
	buf.Write(lenBuf[:])
	buf.Write(key)
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(len(val)))
	buf.Write(lenBuf[:])
	buf.Write(val)
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.file.Write(buf.Bytes()); err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *WAL) Path() string{
	return w.path
}
