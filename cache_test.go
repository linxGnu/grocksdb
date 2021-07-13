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

	cache.DisownData()
}

func TestCacheWithOpts(t *testing.T) {
	opts := NewLRUCacheOptions()
	opts.SetCapacity(19)
	defer opts.Destroy()

	cache := NewLRUCacheWithOptions(opts)
	defer cache.Destroy()

	require.EqualValues(t, 19, cache.GetCapacity())
	cache.SetCapacity(128)
	require.EqualValues(t, 128, cache.GetCapacity())

	cache.DisownData()
}
