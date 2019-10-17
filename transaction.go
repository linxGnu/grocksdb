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

// GetWithCF returns the data associated with the key from the database, with column family, given this transaction.
func (transaction *Transaction) GetWithCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transaction_get_cf(
		transaction.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
	)
	if err = fromCError(cErr); err == nil {
		slice = NewSlice(cValue, cValLen)
	}

	return
}

// GetForUpdate queries the data associated with the key and puts an exclusive lock on the key
// from the database given this transaction.
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

// GetForUpdateWithCF queries the data associated with the key and puts an exclusive lock on the key
// from the database, with column family, given this transaction.
func (transaction *Transaction) GetForUpdateWithCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transaction_get_for_update_cf(
		transaction.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, C.uchar(byte(1)) /*exclusive*/, &cErr,
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

// PutCF writes data associated with a key to the transaction. Key belongs to column family.
func (transaction *Transaction) PutCF(cf *ColumnFamilyHandle, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transaction_put_cf(
		transaction.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// Merge key, value to the transaction.
func (transaction *Transaction) Merge(key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transaction_merge(
		transaction.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// MergeCF key, value to the transaction on specific column family.
func (transaction *Transaction) MergeCF(cf *ColumnFamilyHandle, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transaction_merge_cf(
		transaction.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
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

// DeleteCF removes the data associated with the key (belongs to specific column family) from the transaction.
func (transaction *Transaction) DeleteCF(cf *ColumnFamilyHandle, key []byte) (err error) {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)

	C.rocksdb_transaction_delete_cf(transaction.c, cf.c, cKey, C.size_t(len(key)), &cErr)
	err = fromCError(cErr)

	return
}

// NewIterator returns an iterator that will iterate on all keys in the default
// column family including both keys in the DB and uncommitted keys in this
// transaction.
//
// Setting read_options.snapshot will affect what is read from the
// DB but will NOT change which keys are read from this transaction (the keys
// in this transaction do not yet belong to any snapshot and will be fetched
// regardless).
//
// Caller is responsible for deleting the returned Iterator.
func (transaction *Transaction) NewIterator(opts *ReadOptions) *Iterator {
	return NewNativeIterator(
		unsafe.Pointer(C.rocksdb_transaction_create_iterator(transaction.c, opts.c)))
}

// NewIteratorCF returns an iterator that will iterate on all keys in the specific
// column family including both keys in the DB and uncommitted keys in this
// transaction.
//
// Setting read_options.snapshot will affect what is read from the
// DB but will NOT change which keys are read from this transaction (the keys
// in this transaction do not yet belong to any snapshot and will be fetched
// regardless).
//
// Caller is responsible for deleting the returned Iterator.
func (transaction *Transaction) NewIteratorCF(opts *ReadOptions, cf *ColumnFamilyHandle) *Iterator {
	return NewNativeIterator(
		unsafe.Pointer(C.rocksdb_transaction_create_iterator_cf(transaction.c, opts.c, cf.c)))
}

// SetSavePoint records the state of the transaction for future calls to
// RollbackToSavePoint().  May be called multiple times to set multiple save
// points.
func (transaction *Transaction) SetSavePoint() {
	C.rocksdb_transaction_set_savepoint(transaction.c)
}

// RollbackToSavePoint undo all operations in this transaction (Put, Merge, Delete, PutLogData)
// since the most recent call to SetSavePoint() and removes the most recent
// SetSavePoint().
func (transaction *Transaction) RollbackToSavePoint() (err error) {
	var cErr *C.char
	C.rocksdb_transaction_rollback_to_savepoint(transaction.c, &cErr)
	err = fromCError(cErr)
	return
}

// GetSnapshot returns the Snapshot created by the last call to SetSnapshot().
func (transaction *Transaction) GetSnapshot() *Snapshot {
	return NewNativeSnapshot(C.rocksdb_transaction_get_snapshot(transaction.c))
}

// Destroy deallocates the transaction object.
func (transaction *Transaction) Destroy() {
	C.rocksdb_transaction_destroy(transaction.c)
	transaction.c = nil
}
