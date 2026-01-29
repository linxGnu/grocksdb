package grocksdb_test

import (
	"testing"

	"github.com/linxGnu/grocksdb"
	"github.com/stretchr/testify/require"
)

func TestExportImportFileMetadata(t *testing.T) {
	metadata := grocksdb.NewExportImportFileMetadata()
	defer metadata.Destroy()
	require.NotNil(t, metadata)
}
