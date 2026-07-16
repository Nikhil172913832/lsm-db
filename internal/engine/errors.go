package engine

import "errors"

var InvalidWALFileName = errors.New("Invalid WAL filename")
var InvalidNoOfFilesInWALDir = errors.New("WAL dir contains invalid number of files")
var ErrActiveStateDoubleFree = errors.New("Active state has been freed twice")
