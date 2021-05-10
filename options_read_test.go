package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadOptions(t *testing.T) {
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	require.EqualValues(t, true, ro.VerifyChecksums())
	ro.SetVerifyChecksums(false)
	require.EqualValues(t, false, ro.VerifyChecksums())

	require.EqualValues(t, true, ro.FillCache())
	ro.SetFillCache(false)
	require.EqualValues(t, false, ro.FillCache())

	ro.SetSnapshot(NewNativeSnapshot(nil))

	ro.SetIterateUpperBound([]byte{1, 2, 3})
	ro.SetIterateLowerBound([]byte{1, 1, 1})

	require.EqualValues(t, ReadAllTier, ro.GetReadTier())
	ro.SetReadTier(BlockCacheTier)
	require.EqualValues(t, BlockCacheTier, ro.GetReadTier())

	require.EqualValues(t, false, ro.Tailing())
	ro.SetTailing(true)
	require.EqualValues(t, true, ro.Tailing())

	require.EqualValues(t, 0, ro.GetReadaheadSize())
	ro.SetReadaheadSize(1 << 20)
	require.EqualValues(t, 1<<20, ro.GetReadaheadSize())

	require.EqualValues(t, false, ro.PrefixSameAsStart())
	ro.SetPrefixSameAsStart(true)
	require.EqualValues(t, true, ro.PrefixSameAsStart())

	require.EqualValues(t, false, ro.PinData())
	ro.SetPinData(true)
	require.EqualValues(t, true, ro.PinData())

	require.EqualValues(t, false, ro.GetTotalOrderSeek())
	ro.SetTotalOrderSeek(true)
	require.EqualValues(t, true, ro.GetTotalOrderSeek())

	require.EqualValues(t, 0, ro.GetMaxSkippableInternalKeys())
	ro.SetMaxSkippableInternalKeys(123)
	require.EqualValues(t, 123, ro.GetMaxSkippableInternalKeys())

	require.EqualValues(t, false, ro.GetBackgroundPurgeOnIteratorCleanup())
	ro.SetBackgroundPurgeOnIteratorCleanup(true)
	require.EqualValues(t, true, ro.GetBackgroundPurgeOnIteratorCleanup())

	require.EqualValues(t, false, ro.IgnoreRangeDeletions())
	ro.SetIgnoreRangeDeletions(true)
	require.EqualValues(t, true, ro.IgnoreRangeDeletions())

	require.EqualValues(t, 0, ro.GetDeadline())
	ro.SetDeadline(1000)
	require.EqualValues(t, 1000, ro.GetDeadline())

	require.EqualValues(t, 0, ro.GetIOTimeout())
	ro.SetIOTimeout(1212)
	require.EqualValues(t, 1212, ro.GetIOTimeout())
}
