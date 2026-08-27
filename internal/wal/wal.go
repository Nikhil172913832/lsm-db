package wal

import (
	"encoding/binary"
	"os"
	"sync"

	"github.com/zeebo/xxh3"
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

func (w *WAL) Append(op byte, key, val []byte) error {
	if uint32(len(key)) > MaxKeySize {
		return ErrKeySizeExceedMaxLimit
	}
	if uint32(len(val)) > MaxValueSize {
		return ErrValSizeExceedMaxLimit
	}
	totalSize := 8 + 4 + 1 + 4 + len(key) + len(val)
	buf := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(buf[8:12], uint32(totalSize-12))
	buf[12] = op
	binary.LittleEndian.PutUint32(buf[13:17], uint32(len(key)))
	copy(buf[17:17+len(key)], key)
	copy(buf[17+len(key):], val)
	checksum := xxh3.Hash(buf[12:])
	binary.LittleEndian.PutUint64(buf[0:8], checksum)
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.file.Write(buf); err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *WAL) Path() string {
	return w.path
}
