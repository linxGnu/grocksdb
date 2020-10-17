package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewLRUCache(19)
	defer cache.Destroy()

	require.EqualValues(t, 19, cache.GetCapacity())
	cache.SetCapacity(128)
	require.EqualValues(t, 128, cache.GetCapacity())
}
