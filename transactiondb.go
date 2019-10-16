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

// OpenTransactionDbColumnFamilies opens a database with the specified column families.
func OpenTransactionDbColumnFamilies(
	opts *Options,
	transactionDBOpts *TransactionDBOptions,
	name string,
	cfNames []string,
	cfOpts []*Options,
) (db *TransactionDB, cfHandles []*ColumnFamilyHandle, err error) {
	numColumnFamilies := len(cfNames)
	if numColumnFamilies != len(cfOpts) {
		err = ErrColumnFamilyMustMatch
		return
	}

	cName := C.CString(name)
	cNames := make([]*C.char, numColumnFamilies)
	for i, s := range cfNames {
		cNames[i] = C.CString(s)
	}

	cOpts := make([]*C.rocksdb_options_t, numColumnFamilies)
	for i, o := range cfOpts {
		cOpts[i] = o.c
	}

	cHandles := make([]*C.rocksdb_column_family_handle_t, numColumnFamilies)

	var cErr *C.char
	_db := C.rocksdb_transactiondb_open_column_families(
		opts.c,
		transactionDBOpts.c,
		cName,
		C.int(numColumnFamilies),
		&cNames[0],
		&cOpts[0],
		&cHandles[0],
		&cErr,
	)
	if err = fromCError(cErr); err == nil {
		db = &TransactionDB{
			name:              name,
			c:                 _db,
			opts:              opts,
			transactionDBOpts: transactionDBOpts,
		}
		cfHandles = make([]*ColumnFamilyHandle, numColumnFamilies)
		for i, c := range cHandles {
			cfHandles[i] = NewNativeColumnFamilyHandle(c)
		}
	}

	C.free(unsafe.Pointer(cName))
	for _, s := range cNames {
		C.free(unsafe.Pointer(s))
	}
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

// GetCF returns the data associated with the key from the database, from column family.
func (db *TransactionDB) GetCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (slice *Slice, err error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = byteToChar(key)
	)

	cValue := C.rocksdb_transactiondb_get_cf(
		db.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, &cErr,
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

// Put writes data associated with a key to the database on specific column family.
func (db *TransactionDB) PutCF(opts *WriteOptions, cf *ColumnFamilyHandle, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transactiondb_put_cf(
		db.c, opts.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// Merge writes data associated with a key to the database.
func (db *TransactionDB) Merge(opts *WriteOptions, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transactiondb_merge(
		db.c, opts.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
	)
	err = fromCError(cErr)

	return
}

// MergeCF writes data associated with a key to the database on specific column family.
func (db *TransactionDB) MergeCF(opts *WriteOptions, cf *ColumnFamilyHandle, key, value []byte) (err error) {
	var (
		cErr   *C.char
		cKey   = byteToChar(key)
		cValue = byteToChar(value)
	)

	C.rocksdb_transactiondb_merge_cf(
		db.c, opts.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr,
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

// Delete removes the data associated with the key from the database on specific column family.
func (db *TransactionDB) DeleteCF(opts *WriteOptions, cf *ColumnFamilyHandle, key []byte) (err error) {
	var (
		cErr *C.char
		cKey = byteToChar(key)
	)

	C.rocksdb_transactiondb_delete_cf(db.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cErr)
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

// CreateColumnFamily create a new column family.
func (db *TransactionDB) CreateColumnFamily(opts *Options, name string) (handle *ColumnFamilyHandle, err error) {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)

	cHandle := C.rocksdb_transactiondb_create_column_family(db.c, opts.c, cName, &cErr)
	if err = fromCError(cErr); err == nil {
		handle = NewNativeColumnFamilyHandle(cHandle)
	}

	C.free(unsafe.Pointer(cName))
	return
}

// Write writes a WriteBatch to the database.
func (db *TransactionDB) Write(opts *WriteOptions, batch *WriteBatch) (err error) {
	var cErr *C.char

	C.rocksdb_transactiondb_write(db.c, opts.c, batch.c, &cErr)
	err = fromCError(cErr)

	return
}

// NewIterator returns an Iterator over the the database that uses the
// ReadOptions given.
func (db *TransactionDB) NewIterator(opts *ReadOptions) *Iterator {
	cIter := C.rocksdb_transactiondb_create_iterator(db.c, opts.c)
	return NewNativeIterator(unsafe.Pointer(cIter))
}

// NewIteratorCF returns an Iterator over the the database and column family
// that uses the ReadOptions given.
func (db *TransactionDB) NewIteratorCF(opts *ReadOptions, cf *ColumnFamilyHandle) *Iterator {
	cIter := C.rocksdb_transactiondb_create_iterator_cf(db.c, opts.c, cf.c)
	return NewNativeIterator(unsafe.Pointer(cIter))
}

// Close closes the database.
func (transactionDB *TransactionDB) Close() {
	C.rocksdb_transactiondb_close(transactionDB.c)
	transactionDB.c = nil
}
