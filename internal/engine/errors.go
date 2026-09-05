package engine

import "errors"

var InvalidWALFileName = errors.New("Invalid WAL filename")
var InvalidNoOfFilesInWALDir = errors.New("WAL dir contains invalid number of files")
var ErrActiveStateDoubleFree = errors.New("Active state has been freed twice")
var ErrEmptyKey = errors.New("Key cannot be nil or empty")
var ErrKeySizeExceedMaxLimit = errors.New("key exceeds maximum size")
var ErrValSizeExceedMaxLimit = errors.New("value exceeds maximum size")
