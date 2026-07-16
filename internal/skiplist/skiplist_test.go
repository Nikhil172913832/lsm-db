package skiplist

import (
	"bytes"
	"testing"
)

func TestInsertSingleKey(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte("apple")
	want := []byte("red")
	sl.Insert(key, want)
	node := sl.Search(key)
	if node == nil {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(node.Value, want) {
		t.Fatalf("expected value %q, got %q", want, node.Value)
	}
}

func TestSearchEmptyList(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte("apple")
	node := sl.Search(key)
	if node != nil {
		t.Fatal("expected search to return nil")
	}
}

func TestMissingKey(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte("apple")
	want := []byte("red")
	sl.Insert(key, want)
	keyToSearch := []byte("banana")
	node := sl.Search(keyToSearch)
	if node != nil {
		t.Fatal("expected search to return nil")
	}
}

func TestUpdateExistingKey(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte("apple")
	want := []byte("red")
	sl.Insert(key, want)
	newWant := []byte("green")
	sl.Insert(key, newWant)
	node := sl.Search(key)
	if node == nil {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(node.Value, newWant) {
		t.Fatalf("expected value %q, got %q", newWant, node.Value)
	}
}

func TestInsertMultipleKeys(t *testing.T) {
	sl := NewSkipList(8)
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
		sl.Insert(keys[i], values[i])
	}
	for i := range keys {
		node := sl.Search(keys[i])
		if node == nil {
			t.Fatalf("expected key %q to exist", keys[i])
		}
		if !bytes.Equal(node.Value, values[i]) {
			t.Fatalf(
				"for key %q expected value %q, got %q",
				keys[i],
				values[i],
				node.Value,
			)
		}
	}
}

func TestInsertEmptyKey(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte{}
	want := []byte{}
	sl.Insert(key, want)
	node := sl.Search(key)
	if node == nil {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(node.Value, want) {
		t.Fatalf("expected value %q, got %q", want, node.Value)
	}
}

func TestInsertBinaryKey(t *testing.T) {
	sl := NewSkipList(8)
	key := []byte{0x00, 0xff, 0x10}
	want := []byte{0x10, 0xff, 0x10}
	sl.Insert(key, want)
	node := sl.Search(key)
	if node == nil {
		t.Fatal("expected key to exist")
	}
	if !bytes.Equal(node.Value, want) {
		t.Fatalf("expected value %q, got %q", want, node.Value)
	}
}
