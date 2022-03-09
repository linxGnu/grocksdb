package grocksdb

// #include "rocksdb/c.h"
// #include "grocksdb.h"
import "C"

import (
	"bytes"
)

// Comparing functor.
//
// Three-way comparison. Returns value:
//   < 0 iff "a" < "b",
//   == 0 iff "a" == "b",
//   > 0 iff "a" > "b"
type Comparing = func(a, b []byte) int

// NewComparator creates a Comparator object which contains native c-comparator pointer.
func NewComparator(name string, compare Comparing) *Comparator {
	cmp := &Comparator{name: name, compare: compare}
	idx := registerComperator(cmp)
	cmp.c = C.gorocksdb_comparator_create(C.uintptr_t(idx))
	return cmp
}

// NativeComparator wraps c-comparator pointer.
type Comparator struct {
	c       *C.rocksdb_comparator_t
	compare Comparing
	name    string
}

func (c *Comparator) Compare(a, b []byte) int { return c.compare(a, b) }
func (c *Comparator) Name() string            { return c.name }
func (c *Comparator) Destroy() {
	C.rocksdb_comparator_destroy(c.c)
	c.c = nil
}

// Hold references to comperators.
var comperators = NewCOWList()

type comperatorWrapper struct {
	name       *C.char
	comparator *Comparator
}

func registerComperator(cmp *Comparator) int {
	return comperators.Append(comperatorWrapper{C.CString(cmp.Name()), cmp})
}

//export gorocksdb_comparator_compare
func gorocksdb_comparator_compare(idx int, cKeyA *C.char, cKeyALen C.size_t, cKeyB *C.char, cKeyBLen C.size_t) C.int {
	keyA := charToByte(cKeyA, cKeyALen)
	keyB := charToByte(cKeyB, cKeyBLen)
	return C.int(comperators.Get(idx).(comperatorWrapper).comparator.Compare(keyA, keyB))
}

//export gorocksdb_comparator_name
func gorocksdb_comparator_name(idx int) *C.char {
	return comperators.Get(idx).(comperatorWrapper).name
}

// for testing purpose only
type testBytesReverseComparator struct{}

func (cmp *testBytesReverseComparator) Name() string { return "grocksdb.bytes-reverse" }
func (cmp *testBytesReverseComparator) Compare(a, b []byte) int {
	return bytes.Compare(a, b) * -1
}
func (cmp *testBytesReverseComparator) Destroy() {}
