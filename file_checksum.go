package grocksdb

// #include "rocksdb/c.h"
import "C"

// FileChecksumGenFactory wraps a native file checksum generator factory used to
// compute and verify per-file checksums for SST files.
type FileChecksumGenFactory struct {
	c *C.rocksdb_file_checksum_gen_factory_t
}

// NewCRC32CFileChecksumGenFactory creates a file checksum generator factory
// that produces CRC32C checksums.
func NewCRC32CFileChecksumGenFactory() *FileChecksumGenFactory {
	return &FileChecksumGenFactory{
		c: C.rocksdb_file_checksum_gen_crc32c_factory_create(),
	}
}

// newNativeFileChecksumGenFactory wraps an existing native handle.
func newNativeFileChecksumGenFactory(c *C.rocksdb_file_checksum_gen_factory_t) *FileChecksumGenFactory {
	return &FileChecksumGenFactory{c: c}
}

// Destroy deallocates the FileChecksumGenFactory object.
func (f *FileChecksumGenFactory) Destroy() {
	C.rocksdb_file_checksum_gen_factory_destroy(f.c)
	f.c = nil
}
