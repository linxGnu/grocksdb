package grocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

type ImportColumnFamilyOption struct {
	c *C.rocksdb_import_column_family_options_t
}

func NewImportColumnFamilyOption() *ImportColumnFamilyOption {
	return &ImportColumnFamilyOption{c: C.rocksdb_import_column_family_options_create()}
}

// SetMoveFiles sets to true to move the files instead of copying them.
// The input files will be unlinked after successful ingestion.
// The implementation depends on hard links (LinkFile) instead of traditional
// move (RenameFile) to maximize the chances to restore to the original
// state upon failure.
func (i *ImportColumnFamilyOption) SetMoveFiles(v bool) {
	C.rocksdb_import_column_family_options_set_move_files(i.c, boolToChar(v))
}

func (i *ImportColumnFamilyOption) Destroy() {
	if i.c != nil {
		C.rocksdb_import_column_family_options_destroy(i.c)
		i.c = nil
	}
}
