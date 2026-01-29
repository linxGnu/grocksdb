package grocksdb_test

import (
	"testing"

	"github.com/linxGnu/grocksdb"
)

func TestImportColumnFamilyOption(t *testing.T) {
	o := grocksdb.NewImportColumnFamilyOption()
	defer o.Destroy()

	o.SetMoveFiles(true)
	o.SetMoveFiles(false)
	o.SetMoveFiles(true)
}
