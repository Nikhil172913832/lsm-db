package memtable

import "errors"

var ErrEmptyOrNilValue = errors.New("Cannot insert nil or empty value: nil or empty is reserved for tombstone")
