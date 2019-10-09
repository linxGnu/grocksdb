package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"unsafe"
)

// TransactionDB is a reusable handle to a RocksDB transactional database on disk, created by OpenTransactionDb.
type TransactionDB struct {
	c                 *C.rocksdb_transactiondb_t
	name              string
	opts              *Options
	transactionDBOpts *TransactionDBOptions
}

// OpenTransactionDb opens a database with the specified options.
func OpenTransactionDb(
	opts *Options,
	transactionDBOpts *TransactionDBOptions,
	name string,
) (tdb *TransactionDB, err error) {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)

	db := C.rocksdb_transactiondb_open(
		opts.c, transactionDBOpts.c, cName, &cErr)
	if err = fromCError(cErr); err == nil {
		tdb = &TransactionDB{
			name:              name,
			c:                 db,
			opts:              opts,
			transactionDBOpts: transactionDBOpts,
		}
	}

	C.free(unsafe.Pointer(cName))
	return
}

// NewSnapshot creates a new snapshot of the database.
func (db *TransactionDB) NewSnapshot() *Snapshot {
	return NewNativeSnapshot(C.rocksdb_transactiondb_create_snapshot(db.c))
}

// ReleaseSnapshot releases the snapshot and its resources.
func (db *TransactionDB) ReleaseSnapshot(snapshot *Snapshot) {
	C.rocksdb_transactiondb_release_snapshot(db.c, snapshot.c)
	snapshot.c = nil
}

// TransactionBegin begins a new transaction
// with the WriteOptions and TransactionOptions given.
func (db *TransactionDB) TransactionBegin(
	opts *WriteOptions,
	transactionOpts *TransactionOptions,
	oldTransaction *Transaction,
) *Transaction {
	if oldTransaction != nil {
		return NewNativeTransaction(C.rocksdb_transaction_begin(
			db.c,
			opts.c,
			transactionOpts.c,
			oldTransaction.c,
		))
	}

	return NewNativeTransaction(C.rocksdb_transaction_begin(
		db.c, opts.c, transactionOpts.c, nil))
}

// Get returns the data associated with the key from the database.
func (db *TransactionDB) Get(opts *ReadOptions, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transactiondb_get(
		db.c, opts.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
	)
	if err = fromCError(cErr); err == nil {
		slice = NewSlice(cValue, cValLen)
	}

	return
}

// Put writes data associated with a key to the database.
func (db *TransactionDB) Put(opts *WriteOptions, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transactiondb_put(
		db.c, opts.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// Delete removes the data associated with the key from the database.
func (db *TransactionDB) Delete(opts *WriteOptions, key []byte) (err error) {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)

	C.rocksdb_transactiondb_delete(db.c, opts.c, cKey, C.size_t(len(key)), &cErr)
	err = fromCError(cErr)

	return
}

// NewCheckpoint creates a new Checkpoint for this db.
func (db *TransactionDB) NewCheckpoint() (cp *Checkpoint, err error) {
	var cErr *C.char

	cCheckpoint := C.rocksdb_transactiondb_checkpoint_object_create(
		db.c, &cErr,
	)
	if err = fromCError(cErr); err == nil {
		cp = NewNativeCheckpoint(cCheckpoint)
	}

	return
}

// Close closes the database.
func (transactionDB *TransactionDB) Close() {
	C.rocksdb_transactiondb_close(transactionDB.c)
	transactionDB.c = nil
}
