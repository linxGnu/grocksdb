//go:build !windows
// +build !windows

package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenDbUnix(t *testing.T) {
	db := newTestDB(t, nil)
	defer db.Close()
	require.EqualValues(t, "0", db.GetProperty("rocksdb.num-immutable-mem-table"))
	v, success := db.GetIntProperty("rocksdb.num-immutable-mem-table")
	require.EqualValues(t, uint64(0), v)
	require.True(t, success)
}
