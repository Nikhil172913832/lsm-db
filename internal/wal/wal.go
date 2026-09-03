package wal

import (
	"encoding/binary"
	"io"
	"os"
	"sync"

	"github.com/zeebo/xxh3"
)

type writeRequest struct {
	data     []byte
	err      error
	done     chan struct{}
	isLeader chan struct{}
}
type WAL struct {
	file         *os.File
	path         string
	mu           sync.Mutex
	queue        []*writeRequest
	leaderActive bool
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		file.Close()
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
	data := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(data[8:12], uint32(totalSize-12))
	data[12] = op
	binary.LittleEndian.PutUint32(data[13:17], uint32(len(key)))
	copy(data[17:17+len(key)], key)
	copy(data[17+len(key):], val)
	checksum := xxh3.Hash(data[12:])
	binary.LittleEndian.PutUint64(data[0:8], checksum)
	req := &writeRequest{data: data, done: make(chan struct{}), isLeader: make(chan struct{}, 1)}
	w.mu.Lock()
	w.queue = append(w.queue, req)
	if w.leaderActive {
		w.mu.Unlock()
		select {
		case <-req.done:
			return req.err
		case <-req.isLeader:
			w.mu.Lock()
		}
	} else {
		w.leaderActive = true
	}
	batch := w.queue
	w.queue = nil
	w.mu.Unlock()
	var combined []byte
	for _, r := range batch {
		combined = append(combined, r.data...)
	}
	_, writeErr := w.file.Write(combined)
	syncErr := w.file.Sync()
	var finalErr error
	if writeErr != nil {
		finalErr = writeErr
	} else {
		finalErr = syncErr
	}
	for _, r := range batch {
		r.err = finalErr
		close(r.done)
	}
	w.mu.Lock()
	if len(w.queue) == 0 {
		w.leaderActive = false
	} else {
		nextLeader := w.queue[0]
		nextLeader.isLeader <- struct{}{}
	}
	w.mu.Unlock()
	return req.err
}

func (w *WAL) Path() string {
	return w.path
}
