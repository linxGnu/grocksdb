package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import "unsafe"

// LiveFile live table (SST) file.
type LiveFile struct {
	c *C.rocksdb_livefile_t
}

func NewLiveFile() *LiveFile {
	return &LiveFile{
		c: C.rocksdb_livefile_create(),
	}
}

func (l *LiveFile) SetColumnFamilyName(name string) {
	cName := C.CString(name)
	C.rocksdb_livefile_set_column_family_name(l.c, cName)
	C.free(unsafe.Pointer(cName))
}

func (l *LiveFile) SetLevel(level int) {
	C.rocksdb_livefile_set_level(l.c, C.int(level))
}

func (l *LiveFile) SetName(name string) {
	cName := C.CString(name)
	C.rocksdb_livefile_set_name(l.c, cName)
	C.free(unsafe.Pointer(cName))
}

func (l *LiveFile) SetDirectory(dir string) {
	cName := C.CString(dir)
	C.rocksdb_livefile_set_directory(l.c, cName)
	C.free(unsafe.Pointer(cName))
}

func (l *LiveFile) SetSize(size int) {
	C.rocksdb_livefile_set_size(l.c, C.size_t(size))
}

func (l *LiveFile) SetSmallestKey(key []byte) {
	k := refGoBytes(key)
	C.rocksdb_livefile_set_smallest_key(l.c, k, C.size_t(len(key)))
}

func (l *LiveFile) SetLargestKey(key []byte) {
	k := refGoBytes(key)
	C.rocksdb_livefile_set_largest_key(l.c, k, C.size_t(len(key)))
}

func (l *LiveFile) SetSmallestSeqNo(seq uint64) {
	C.rocksdb_livefile_set_smallest_seqno(l.c, C.uint64_t(seq))
}

func (l *LiveFile) SetLargestSeqNo(seq uint64) {
	C.rocksdb_livefile_set_largest_seqno(l.c, C.uint64_t(seq))
}

func (l *LiveFile) SetNumEntries(n uint64) {
	C.rocksdb_livefile_set_num_entries(l.c, C.uint64_t(n))
}

func (l *LiveFile) SetNumDeletions(n uint64) {
	C.rocksdb_livefile_set_num_deletions(l.c, C.uint64_t(n))
}

func (l *LiveFile) Destroy() {
	if l.c != nil {
		C.rocksdb_livefile_destroy(l.c)
		l.c = nil
	}
}

// LiveFiles is a list of all live table (SST) files and how they fit into the
// LSM-trees, such as column family, level, key range, etc.
type LiveFiles struct {
	c *C.rocksdb_livefiles_t
}

func NewLiveFiles() *LiveFiles {
	return NewNativeLiveFiles(C.rocksdb_livefiles_create())
}

func NewNativeLiveFiles(c *C.rocksdb_livefiles_t) *LiveFiles {
	return &LiveFiles{c: c}
}

// Count returns number of live table (SST) files.
func (l *LiveFiles) Count() int {
	return int(C.rocksdb_livefiles_count(l.c))
}

// ColumnFamilyName returns name of the column family of a live file.
func (l *LiveFiles) ColumnFamilyName(liveFileIndex int) string {
	// returning const char* -> do not C.free
	cName := C.rocksdb_livefiles_column_family_name(l.c, C.int(liveFileIndex))
	return C.GoString(cName)
}

// FileName returns name of a live file.
func (l *LiveFiles) Name(liveFileIndex int) string {
	cName := C.rocksdb_livefiles_name(l.c, C.int(liveFileIndex))
	return C.GoString(cName)
}

// Directory returns the directory which contains the live file.
func (l *LiveFiles) Directory(liveFileIndex int) string {
	cName := C.rocksdb_livefiles_directory(l.c, C.int(liveFileIndex))
	return C.GoString(cName)
}

// Level returns the level at which the live file resides.
func (l *LiveFiles) Level(liveFileIndex int) int {
	return int(C.rocksdb_livefiles_level(l.c, C.int(liveFileIndex)))
}

// Size returns the size of the live file.
func (l *LiveFiles) Size(liveFileIndex int) int {
	return int(C.rocksdb_livefiles_size(l.c, C.int(liveFileIndex)))
}

// SmallestKey returns the smallest key in the live file.
func (l *LiveFiles) SmallestKey(liveFileIndex int) []byte {
	var cValLen C.size_t

	// returning const char* -> do not C.free
	cValue := C.rocksdb_livefiles_smallestkey(l.c, C.int(liveFileIndex), &cValLen)

	return C.GoBytes(unsafe.Pointer(&cValue), C.int(cValLen))
}

// LargestKey returns the largest key in the live file.
func (l *LiveFiles) LargestKey(liveFileIndex int) []byte {
	var cValLen C.size_t

	// returning const char* -> do not C.free
	cValue := C.rocksdb_livefiles_largestkey(l.c, C.int(liveFileIndex), &cValLen)

	return C.GoBytes(unsafe.Pointer(&cValue), C.int(cValLen))
}

// SmallestSeqNo returns smallest sequence number in the live file.
func (l *LiveFiles) SmallestSeqNo(liveFileIndex int) uint64 {
	return uint64(C.rocksdb_livefiles_smallest_seqno(l.c, C.int(liveFileIndex)))
}

// LargestSeqNo returns largest sequence number in the live file.
func (l *LiveFiles) LargestSeqNo(liveFileIndex int) uint64 {
	return uint64(C.rocksdb_livefiles_largest_seqno(l.c, C.int(liveFileIndex)))
}

// NumEntries returns number of entries in the live file.
func (l *LiveFiles) NumEntries(liveFileIndex int) uint64 {
	return uint64(C.rocksdb_livefiles_entries(l.c, C.int(liveFileIndex)))
}

// NumDeletions returns number of deletions in the live file.
func (l *LiveFiles) NumDeletions(liveFileIndex int) uint64 {
	return uint64(C.rocksdb_livefiles_deletions(l.c, C.int(liveFileIndex)))
}

// AddLiveFile to the collection.
func (l *LiveFiles) AddLiveFile(lf *LiveFile) {
	C.rocksdb_livefiles_add(l.c, lf.c)
}

// Destroy LiveFiles.
func (l *LiveFiles) Destroy() {
	if l.c != nil {
		C.rocksdb_livefiles_destroy(l.c)
		l.c = nil
	}
}
