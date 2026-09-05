package engine

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestEngineBasicPutGetDelete(t *testing.T) {
	tempDir := t.TempDir()
	walDir := filepath.Join(tempDir, "wal")
	sstableDir := filepath.Join(tempDir, "sstable")
	os.MkdirAll(walDir, 0755)
	os.MkdirAll(sstableDir, 0755)
	e, err := NewEngine(walDir, sstableDir, 8, 10, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	err = e.Put([]byte("fruit"), []byte("apple"))
	if err != nil {
		t.Fatal(err)
	}
	val, found, err := e.Get([]byte("fruit"))
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Fatal("Expected key 'fruit' to be found")
	}
	expected := []byte("apple")
	if !bytes.Equal(val, expected) {
		t.Fatalf("Expected %q, got %q", expected, val)
	}
	err = e.Put([]byte("emptyVal"), []byte{})
	if err != nil {
		t.Fatal(err)
	}
	val, found, err = e.Get([]byte("emptyVal"))
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Fatal("Expected key 'emptyVal' to be found")
	}
	expected = []byte{}
	if !bytes.Equal(val, expected) {
		t.Fatalf("Expected %q, got %q", expected, val)
	}
	err = e.Delete([]byte("fruit"))
	if err != nil {
		t.Fatal(err)
	}
	val, found, err = e.Get([]byte("fruit"))
	if err != nil {
		t.Fatalf("Get failed after delete: %v", err)
	}
	if found {
		t.Fatalf("Expected key 'fruit' to not exist, got value %q", val)
	}
	if val != nil {
		t.Fatalf("Expected nil value for deleted key, got %q", val)
	}
	val, found, err = e.Get([]byte("unknown"))
	if err != nil {
		t.Fatalf("Get failed for unknown key: %v", err)
	}
	if found {
		t.Fatalf("Expected unknown key to not be found, got value %q", val)
	}
	if val != nil {
		t.Fatalf("Expected nil value for unknown key, got %q", val)
	}
	err = e.Put(nil, []byte("val"))
	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("Expected nil key to fail with ErrEmptyKey, got %v", err)
	}
	err = e.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEngineMemtableRotation(t *testing.T) {
	tempDir := t.TempDir()
	walDir := filepath.Join(tempDir, "wal")
	sstableDir := filepath.Join(tempDir, "sstable")
	os.MkdirAll(walDir, 0755)
	os.MkdirAll(sstableDir, 0755)
	e, err := NewEngine(walDir, sstableDir, 8, 100, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 30 {
		e.Put(fmt.Appendf(nil, "key%d", i), fmt.Appendf(nil, "value%d", i))
	}
	if len(e.state.immutable) == 0 {
		t.Fatal("No rotation took place")
	}
	err = e.Delete([]byte("key0"))
	if err != nil {
		t.Fatal(err)
	}
	if _, found, _ := e.Get([]byte("key0")); found {
		t.Fatal("Expected key to be deleted")
	}
	for i := 1; i < 30; i++ {
		val, found, err := e.Get(fmt.Appendf(nil, "key%d", i))
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if !found {
			t.Fatalf("Expected key%d to be found", i)
		}
		expected := fmt.Appendf(nil, "value%d", i)
		if !bytes.Equal(val, expected) {
			t.Fatalf("Expected %q, got %q", expected, val)
		}
	}
}

func TestEngine_CrashRecovery(t *testing.T) {
	tempDir := t.TempDir()
	walDir := filepath.Join(tempDir, "wal")
	sstableDir := filepath.Join(tempDir, "sstable")
	if err := os.MkdirAll(walDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(sstableDir, 0755); err != nil {
		t.Fatal(err)
	}
	e, err := NewEngine(walDir, sstableDir, 8, 10, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 15; i++ {
		key := fmt.Appendf(nil, "key%d", i)
		value := fmt.Appendf(nil, "value%d", i)
		if err := e.Put(key, value); err != nil {
			t.Fatalf("Put failed for key%d: %v", i, err)
		}
	}
	for i := 1; i <= 3; i++ {
		key := fmt.Appendf(nil, "key%d", i)
		if err := e.Delete(key); err != nil {
			t.Fatalf("Delete failed for key%d: %v", i, err)
		}
	}
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	e, err = NewEngine(walDir, sstableDir, 8, 10, 1<<20, 64<<20)
	if err != nil {
		t.Fatalf("Failed to reopen engine: %v", err)
	}
	defer e.Close()
	for i := 1; i <= 3; i++ {
		key := fmt.Appendf(nil, "key%d", i)
		val, found, err := e.Get(key)
		if err != nil {
			t.Fatalf("Get failed for deleted key%d: %v", i, err)
		}
		if found {
			t.Fatalf("Expected deleted key%d to not be found, got %q", i, val)
		}
		if val != nil {
			t.Fatalf("Expected nil value for deleted key%d, got %q", i, val)
		}
	}
	for i := 4; i <= 15; i++ {
		key := fmt.Appendf(nil, "key%d", i)
		expected := fmt.Appendf(nil, "value%d", i)
		val, found, err := e.Get(key)
		if err != nil {
			t.Fatalf("Get failed for key%d: %v", i, err)
		}
		if !found {
			t.Fatalf("Expected key%d to be found after recovery", i)
		}
		if !bytes.Equal(val, expected) {
			t.Fatalf("Expected key%d to have value %q, got %q", i, expected, val)
		}
	}
}

func TestEngine_ConcurrentReadWrite(t *testing.T) {
	tempDir := t.TempDir()
	walDir := filepath.Join(tempDir, "wal")
	sstableDir := filepath.Join(tempDir, "sstable")
	if err := os.MkdirAll(walDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(sstableDir, 0755); err != nil {
		t.Fatal(err)
	}
	e, err := NewEngine(walDir, sstableDir, 8, 10, 1<<20, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	const numGoroutines = 20
	const operations = 1000
	const numKeys = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < operations; i++ {
				keyID := (id + i) % numKeys
				key := fmt.Appendf(nil, "key%d", keyID)
				value := fmt.Appendf(nil, "value-%d-%d", id, i)
				switch i % 3 {
				case 0:
					if err := e.Put(key, value); err != nil {
						t.Errorf("goroutine %d: Put failed: %v", id, err)
					}
				case 1:
					if _, _, err := e.Get(key); err != nil {
						t.Errorf("goroutine %d: Get failed: %v", id, err)
					}
				case 2:
					if err := e.Delete(key); err != nil {
						t.Errorf("goroutine %d: Delete failed: %v", id, err)
					}
				}
			}
		}(g)
	}
	wg.Wait()
}
