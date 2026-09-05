package wal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/Nikhil172913832/lsm-db/internal/base"
)

func TestAppendAndReadAll(t *testing.T) {
	w, err := NewWAL(filepath.Join(t.TempDir(), "wal"), 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err.Error())
	}
	keys := [][]byte{
		[]byte("mango"),
		[]byte("app"),
		[]byte("orange"),
		[]byte("banana"),
		[]byte("application"),
		[]byte("kiwi"),
		[]byte("l"),
		[]byte("pears"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
		[]byte("green"),
		[]byte("llllll"),
		[]byte("Soap"),
	}
	for i := range keys {
		w.Append(base.OpPut, keys[i], values[i])
	}
	type record struct {
		op    byte
		key   []byte
		value []byte
	}
	var collected []record
	err = w.ReadAll(func(op byte, key, value []byte) error {
		collected = append(collected, record{
			op:    op,
			key:   key,
			value: value,
		})
		return nil
	})
	for i := range keys {
		if collected[i].op != base.OpPut {
			t.Fatal("Incorrect op type")
		}
		if !bytes.Equal(collected[i].key, keys[i]) {
			t.Fatal("Incorrect key")
		}
		if !bytes.Equal(collected[i].value, values[i]) {
			t.Fatal("Incorrect Value")
		}
	}
}

func TestBitFlipCorruption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal")
	w, err := NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err.Error())
	}
	keys := [][]byte{
		[]byte("mango"),
		[]byte("app"),
		[]byte("orange"),
		[]byte("banana"),
		[]byte("application"),
		[]byte("kiwi"),
		[]byte("l"),
		[]byte("pears"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
		[]byte("green"),
		[]byte("llllll"),
		[]byte("Soap"),
	}
	for i := range keys {
		w.Append(base.OpPut, keys[i], values[i])
	}
	type record struct {
		op    byte
		key   []byte
		value []byte
	}
	var collected []record
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	var b [1]byte
	file.ReadAt(b[:], 15)
	b[0] ^= 0xFF
	file.WriteAt(b[:], 15)
	file.Sync()
	file.Close()
	err = w.ReadAll(func(op byte, key, value []byte) error {
		collected = append(collected, record{
			op:    op,
			key:   key,
			value: value,
		})
		return nil
	})
	if !errors.Is(err, ErrCorruptRecord) {
		t.Fatal("Expected ReadAll to fail")
	}
}

func TestTornWriteTruncation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal")
	w, err := NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	keys := [][]byte{
		[]byte("mango"),
		[]byte("app"),
		[]byte("orange"),
		[]byte("banana"),
		[]byte("application"),
		[]byte("kiwi"),
		[]byte("l"),
		[]byte("pears"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
		[]byte("green"),
		[]byte("llllll"),
		[]byte("Soap"),
	}
	for i := range keys {
		w.Append(base.OpPut, keys[i], values[i])
	}
	type record struct {
		op    byte
		key   []byte
		value []byte
	}
	var collected []record
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fileSize := fileInfo.Size()
	file.Truncate(fileSize - 5)
	file.Sync()
	file.Close()
	w, err = NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	err = w.ReadAll(func(op byte, key, value []byte) error {
		collected = append(collected, record{
			op:    op,
			key:   key,
			value: value,
		})
		return nil
	})
	if !errors.Is(err, nil) {
		t.Fatal("Expected ReadAll to fail")
	}
	for i := range 7 {
		if collected[i].op != base.OpPut {
			t.Fatal("Incorrect op type")
		}
		if !bytes.Equal(collected[i].key, keys[i]) {
			t.Fatal("Incorrect key")
		}
		if !bytes.Equal(collected[i].value, values[i]) {
			t.Fatal("Incorrect Value")
		}
	}
}

func TestReopenAndAppendContinuity(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal")
	w, err := NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err.Error())
	}
	keys := [][]byte{
		[]byte("mango"),
		[]byte("app"),
		[]byte("orange"),
		[]byte("banana"),
		[]byte("application"),
		[]byte("kiwi"),
		[]byte("l"),
		[]byte("pears"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
		[]byte("green"),
		[]byte("llllll"),
		[]byte("Soap"),
	}
	for i := range 4 {
		w.Append(base.OpPut, keys[i], values[i])
	}
	w.file.Close()
	w, err = NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 4; i < 8; i++ {
		w.Append(base.OpPut, keys[i], values[i])
	}
	type record struct {
		op    byte
		key   []byte
		value []byte
	}
	var collected []record
	err = w.ReadAll(func(op byte, key, value []byte) error {
		collected = append(collected, record{
			op:    op,
			key:   key,
			value: value,
		})
		return nil
	})
	for i := range keys {
		if collected[i].op != base.OpPut {
			t.Fatal("Incorrect op type")
		}
		if !bytes.Equal(collected[i].key, keys[i]) {
			t.Fatal("Incorrect key")
		}
		if !bytes.Equal(collected[i].value, values[i]) {
			t.Fatal("Incorrect Value")
		}
	}
}

func TestConcurrentGroupCommit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal_concurrent")
	w, err := NewWAL(path, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	const numWriters = 20
	const writesPerGoroutine = 50
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range numWriters {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			<-start
			for j := 0; j < writesPerGoroutine; j++ {
				key := fmt.Appendf(nil, "w_%02d_k_%04d", writerID, j)
				val := fmt.Appendf(nil, "val_%02d_%04d", writerID, j)
				if err := w.Append(base.OpPut, key, val); err != nil {
					t.Errorf("Append failed: %v", err)
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
	count := 0
	err = w.ReadAll(func(op byte, key, value []byte) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if count != numWriters*writesPerGoroutine {
		t.Fatalf("Expected %d records, got %d", numWriters*writesPerGoroutine, count)
	}
}
