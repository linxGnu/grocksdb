package grocksdb

// #include "rocksdb/c.h"
import "C"

// ColumnFamilyMetadata contains metadata info of column family.
type ColumnFamilyMetadata struct {
	c *C.rocksdb_column_family_metadata_t
}

// Destroy releases allocated memory for this instance.
func (cm *ColumnFamilyMetadata) Destroy() {
	if cm.c != nil {
		C.rocksdb_column_family_metadata_destroy(cm.c)
		cm.c = nil
	}
}
