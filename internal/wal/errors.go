package wal

import "errors"

var ErrKeySizeExceedMaxLimit = errors.New("key exceeds maximum size")
var ErrValSizeExceedMaxLimit = errors.New("value exceeds maximum size")
var ErrCorruptRecord = errors.New("Corrupt WAL record")
