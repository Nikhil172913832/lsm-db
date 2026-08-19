package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
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

func TestIterator_EmptyList(t *testing.T) {
	sl := NewSkipList(8)
	it := sl.NewIterator()
	key := []byte("apple")
	it.Seek(key)
	if it.Valid() {
		t.Fatal("Expected valid to be false")
	}
	it.SeekToFirst()
	if it.Valid() {
		t.Fatal("Expected valid to be false")
	}
	it.SeekToLast()
	if it.Valid() {
		t.Fatal("Expected valid to be false")
	}
}

func TestIterator_FullForwardScan(t *testing.T) {
	sl := NewSkipList(8)
	it := sl.NewIterator()
	keys := [][]byte{
		[]byte("banana"),
		[]byte("apple"),
		[]byte("date"),
		[]byte("cherry"),
		[]byte("eldenberry"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
	}
	for i := range keys {
		sl.Insert(keys[i], values[i])
	}
	it.SeekToFirst()
	prev := it.Key()
	it.Next()
	for {
		if it.Valid() {
			current := it.Key()
			if bytes.Compare(prev, current) < 0 {
				prev = current
				it.Next()
			} else {
				t.Fatal("Keys are not sorted")
			}
		} else {
			break
		}
	}
}

func TestIterator_SeekToLast(t *testing.T) {
	sl := NewSkipList(8)
	it := sl.NewIterator()
	keys := [][]byte{
		[]byte("banana"),
		[]byte("apple"),
		[]byte("date"),
		[]byte("cherry"),
		[]byte("eldenberry"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
	}
	for i := range keys {
		sl.Insert(keys[i], values[i])
	}
	it.SeekToLast()
	if !bytes.Equal(it.Key(), keys[4]) {
		t.Fatalf("Expected %q key on SeekToLast got %q", keys[4], it.Key())
	}
}

func TestIterator_Seek(t *testing.T) {
	sl := NewSkipList(8)
	it := sl.NewIterator()
	keys := [][]byte{
		[]byte("10"),
		[]byte("20"),
		[]byte("30"),
		[]byte("40"),
		[]byte("50"),
	}
	values := [][]byte{
		[]byte("yellow"),
		[]byte("MyApp"),
		[]byte("MyOrange"),
		[]byte("banananana"),
		[]byte("MyApplication"),
	}
	for i := range keys {
		sl.Insert(keys[i], values[i])
	}
	it.Seek(keys[1])
	if !bytes.Equal(it.Key(), keys[1]) {
		t.Fatalf("Expected seek to return %q got %q", keys[1], it.Key())
	}
	it.Seek([]byte("25"))
	if !bytes.Equal(it.Key(), keys[2]) {
		t.Fatalf("Expected seek to return %q got %q", keys[2], it.Key())
	}
	it.Seek([]byte("05"))
	if !bytes.Equal(it.Key(), keys[0]) {
		t.Fatalf("Expected seek to return %q got %q", keys[0], it.Key())
	}
	it.Seek([]byte("55"))
	if it.Valid() {
		t.Fatal("Expected valid to return false")
	}
}

func TestIterator_LargeMonotonicScan(t *testing.T) {
	sl := NewSkipList(8)
	it := sl.NewIterator()
	keys := rand.Perm(1000)
	values := rand.Perm(1000)
	for i := range keys{
		sl.Insert(fmt.Appendf(nil, "%03d", keys[i]), fmt.Appendf(nil, "%03d", values[i]))
	}
	it.SeekToFirst()
	prev := it.Key()
	it.Next()
	for {
		if it.Valid() {
			current := it.Key()
			if bytes.Compare(prev, current) < 0 {
				prev = current
				it.Next()
			} else {
				t.Fatal("Keys are not sorted")
			}
		} else {
			break
		}
	}
}
