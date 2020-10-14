package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

// ColumnFamilyHandle represents a handle to a ColumnFamily.
type ColumnFamilyHandle struct {
	c *C.rocksdb_column_family_handle_t
}

// NewNativeColumnFamilyHandle creates a ColumnFamilyHandle object.
func NewNativeColumnFamilyHandle(c *C.rocksdb_column_family_handle_t) *ColumnFamilyHandle {
	return &ColumnFamilyHandle{c}
}

// Destroy calls the destructor of the underlying column family handle.
func (h *ColumnFamilyHandle) Destroy() {
	C.rocksdb_column_family_handle_destroy(h.c)
	h.c = nil
}

// ColumnFamilyHandles represents collection of multiple column family handle.
type ColumnFamilyHandles []*ColumnFamilyHandle

func (cfs ColumnFamilyHandles) toCSlice() columnFamilySlice {
	cCFs := make(columnFamilySlice, len(cfs))
	for i, cf := range cfs {
		cCFs[i] = cf.c
	}
	return cCFs
}
