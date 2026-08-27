package wal

import (
	"encoding/binary"

	"github.com/zeebo/xxh3"
)

func (w *WAL) readRecord(offset int64) (op byte, key []byte, val []byte, nextOffset int64, err error) {
	start := offset
	cursor := offset
	var header [12]byte
	if _, err := w.file.ReadAt(header[:], cursor); err != nil {
		return 0, nil, nil, start, err
	}
	checksum := binary.LittleEndian.Uint64(header[0:8])
	payloadLen := binary.LittleEndian.Uint32(header[8:12])
	if payloadLen > MaxPayloadSize {
		return 0, nil, nil, start, ErrCorruptRecord
	}
	cursor += 12
	payload := make([]byte, payloadLen)
	if _, err := w.file.ReadAt(payload, cursor); err != nil {
		return 0, nil, nil, start, err
	}
	nextOffset = cursor + int64(payloadLen)
	actualChecksum := xxh3.Hash(payload)
	if actualChecksum != checksum {
		return 0, nil, nil, start, ErrCorruptRecord
	}
	op = payload[0]
	keyLen := binary.LittleEndian.Uint32(payload[1:5])
	key = payload[5 : 5+keyLen]
	val = payload[5+keyLen:]
	return
}
