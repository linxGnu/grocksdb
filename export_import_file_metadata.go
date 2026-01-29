package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"unsafe"
)

// ExportImportFileMetadata is the metadata returned as output from ExportColumnFamily()
// and used as input to CreateColumnFamiliesWithImport().
type ExportImportFileMetadata struct {
	c *C.rocksdb_export_import_files_metadata_t
}

func NewExportImportFileMetadata() *ExportImportFileMetadata {
	return NewNativeExportImportFileMetadata(C.rocksdb_export_import_files_metadata_create())
}

func NewNativeExportImportFileMetadata(c *C.rocksdb_export_import_files_metadata_t) *ExportImportFileMetadata {
	return &ExportImportFileMetadata{c: c}
}

func (e *ExportImportFileMetadata) GetComparatorName() string {
	cValue := C.rocksdb_export_import_files_metadata_get_db_comparator_name(e.c)
	name := C.GoString(cValue)
	C.free(unsafe.Pointer(cValue))

	return name
}

func (e *ExportImportFileMetadata) SetComparatorName(name string) {
	cName := C.CString(name)
	C.rocksdb_export_import_files_metadata_set_db_comparator_name(e.c, cName)
	C.free(unsafe.Pointer(cName))
}

// GetFiles obtains a list of all live table (SST) files and how they fit into the
// LSM-trees, such as column family, level, key range, etc.
func (e *ExportImportFileMetadata) GetFiles() *LiveFiles {
	return NewNativeLiveFiles(C.rocksdb_export_import_files_metadata_get_files(e.c))
}

// SetFiles set live files.
func (e *ExportImportFileMetadata) SetFiles(files *LiveFiles) {
	C.rocksdb_export_import_files_metadata_set_files(e.c, files.c)
}

// Destroy ExportImportFileMetadata.
func (e *ExportImportFileMetadata) Destroy() {
	if e.c != nil {
		C.rocksdb_export_import_files_metadata_destroy(e.c)
		e.c = nil
	}
}
