package grocksdb

// #include "rocksdb/c.h"
import "C"

// SSTPartitionerFactory wraps a native SST partitioner factory. An SST
// partitioner decides where output files should be cut during compaction.
type SSTPartitionerFactory struct {
	c *C.rocksdb_sst_partitioner_factory_t
}

// NewFixedPrefixSSTPartitionerFactory creates an SST partitioner factory that
// forces SST files to be cut on a fixed-length key prefix boundary, so that
// keys sharing the same prefix never span two files.
func NewFixedPrefixSSTPartitionerFactory(prefixLen uint) *SSTPartitionerFactory {
	return &SSTPartitionerFactory{
		c: C.rocksdb_sst_partitioner_fixed_prefix_factory_create(C.size_t(prefixLen)),
	}
}

// newNativeSSTPartitionerFactory wraps an existing native handle.
func newNativeSSTPartitionerFactory(c *C.rocksdb_sst_partitioner_factory_t) *SSTPartitionerFactory {
	return &SSTPartitionerFactory{c: c}
}

// Destroy deallocates the SSTPartitionerFactory object.
func (f *SSTPartitionerFactory) Destroy() {
	C.rocksdb_sst_partitioner_factory_destroy(f.c)
	f.c = nil
}
