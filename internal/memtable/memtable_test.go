package memtable

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
)

func TestInsertSingleKey(t *testing.T) {
	mt := NewMemtable(8, 10)
	key := []byte("apple")
	want := []byte("red")
	mt.Put(key, want)
	got, found := mt.Get(key)
	if !found {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected value %q, got %q", want, got)
	}
}

func TestInsertEmptySlice(t *testing.T) {
	mt := NewMemtable(8, 10)
	key := []byte("apple")
	want := []byte{}
	mt.Put(key, want)
	got, found := mt.Get(key)
	if !found {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("expected value %q, got %q", want, got)
	}
}

func TestDeleteAKey(t *testing.T) {
	mt := NewMemtable(8, 10)
	key := []byte("apple")
	want := []byte("red")
	mt.Put(key, want)
	mt.Delete(key)
	got, found := mt.Get(key)
	if got != nil || !found {
		t.Fatalf("expected nil value got %q", got)
	}
}

func TestGetMissingKey(t *testing.T) {
	mt := NewMemtable(8, 10)
	key := []byte("apple")
	got, found := mt.Get(key)
	if found {
		t.Fatalf("expected key to not exist got %q", got)
	}
}

func TestConcurrency(t *testing.T) {
	mt := NewMemtable(8, 100)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range 20 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			key := []byte(strconv.Itoa(i))
			value := []byte("Ding Dong")
			mt.Put(key, value)

		}(i)
	}
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			key := []byte(strconv.Itoa(i))
			mt.Get(key)
			
		}()
	}
	close(start)
	wg.Wait()
}
