package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"unsafe"
)

// Transaction is used with TransactionDB for transaction support.
type Transaction struct {
	c *C.rocksdb_transaction_t
}

// NewNativeTransaction creates a Transaction object.
func NewNativeTransaction(c *C.rocksdb_transaction_t) *Transaction {
	return &Transaction{c}
}

// Commit commits the transaction to the database.
func (transaction *Transaction) Commit() (err error) {
	var cErr *C.char
	C.rocksdb_transaction_commit(transaction.c, &cErr)
	err = fromCError(cErr)
	return
}

// Rollback performs a rollback on the transaction.
func (transaction *Transaction) Rollback() (err error) {
	var cErr *C.char
	C.rocksdb_transaction_rollback(transaction.c, &cErr)
	err = fromCError(cErr)
	return
}

// Get returns the data associated with the key from the database given this transaction.
func (transaction *Transaction) Get(opts *ReadOptions, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transaction_get(
		transaction.c, opts.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
	)
	if err = fromCError(cErr); err == nil {
		slice = NewSlice(cValue, cValLen)
	}

	return
}

// GetForUpdate queries the data associated with the key and puts an exclusive lock on the key from the database given this transaction.
func (transaction *Transaction) GetForUpdate(opts *ReadOptions, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transaction_get_for_update(
		transaction.c, opts.c, cKey, C.size_t(len(key)), &cValLen, C.uchar(byte(1)) /*exclusive*/, &cErr,
	)
	if err = fromCError(cErr); err == nil {
		slice = NewSlice(cValue, cValLen)
	}

	return
}

// Put writes data associated with a key to the transaction.
func (transaction *Transaction) Put(key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transaction_put(
		transaction.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// Delete removes the data associated with the key from the transaction.
func (transaction *Transaction) Delete(key []byte) (err error) {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)

	C.rocksdb_transaction_delete(transaction.c, cKey, C.size_t(len(key)), &cErr)
	err = fromCError(cErr)

	return
}

// NewIterator returns an Iterator over the database that uses the
// ReadOptions given.
func (transaction *Transaction) NewIterator(opts *ReadOptions) *Iterator {
	return NewNativeIterator(
		unsafe.Pointer(C.rocksdb_transaction_create_iterator(transaction.c, opts.c)))
}

// Destroy deallocates the transaction object.
func (transaction *Transaction) Destroy() {
	C.rocksdb_transaction_destroy(transaction.c)
	transaction.c = nil
}
