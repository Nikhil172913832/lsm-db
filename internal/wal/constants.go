package wal

const (
	MaxKeySize   = 1 << 20
	MaxValueSize = 64 << 20
	MaxPayloadSize = MaxKeySize+MaxValueSize
	OpPut    byte = 1
    OpDelete byte = 2
)
