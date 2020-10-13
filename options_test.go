package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	opts := NewDefaultOptions()
	defer opts.Destroy()

	cto := NewCuckooTableOptions()
	opts.SetCuckooTableFactory(cto)
	cto.Destroy()

	opts.SetDumpMallocStats(true)
	opts.SetMemtableWholeKeyFiltering(true)

	opts.SetMaxBackgroundJobs(10)
	require.EqualValues(t, 10, opts.GetMaxBackgroundJobs())

	opts.SetMaxBackgroundCompactions(9)
	require.EqualValues(t, 9, opts.GetMaxBackgroundCompactions())

	opts.SetBaseBackgroundCompactions(4)
	require.EqualValues(t, 4, opts.GetBaseBackgroundCompactions())

	opts.SetMaxBackgroundFlushes(8)
	require.EqualValues(t, 8, opts.GetMaxBackgroundFlushes())

	opts.SetMaxLogFileSize(1 << 30)
	require.EqualValues(t, 1<<30, opts.GetMaxLogFileSize())

	opts.SetLogFileTimeToRoll(924)
	require.EqualValues(t, 924, opts.GetLogFileTimeToRoll())

	opts.SetKeepLogFileNum(19)
	require.EqualValues(t, 19, opts.GetKeepLogFileNum())

	opts.SetRecycleLogFileNum(81)
	require.EqualValues(t, 81, opts.GetRecycleLogFileNum())

	opts.SetSoftRateLimit(0.8)
	require.EqualValues(t, 0.8, opts.GetSoftRateLimit())

	opts.SetHardRateLimit(0.5)
	require.EqualValues(t, 0.5, opts.GetHardRateLimit())

	opts.SetSoftPendingCompactionBytesLimit(50 << 18)
	require.EqualValues(t, 50<<18, opts.GetSoftPendingCompactionBytesLimit())

	opts.SetHardPendingCompactionBytesLimit(50 << 19)
	require.EqualValues(t, 50<<19, opts.GetHardPendingCompactionBytesLimit())

	opts.SetRateLimitDelayMaxMilliseconds(5000)
	require.EqualValues(t, 5000, opts.GetRateLimitDelayMaxMilliseconds())

	opts.SetMaxManifestFileSize(23 << 10)
	require.EqualValues(t, 23<<10, opts.GetMaxManifestFileSize())

	opts.SetTableCacheNumshardbits(5)
	require.EqualValues(t, 5, opts.GetTableCacheNumshardbits())

	opts.SetArenaBlockSize(9 << 20)
	require.EqualValues(t, 9<<20, opts.GetArenaBlockSize())

	opts.SetUseFsync(true)
	require.EqualValues(t, true, opts.UseFsync())

	opts.SetLevelCompactionDynamicLevelBytes(true)
	require.EqualValues(t, true, opts.GetLevelCompactionDynamicLevelBytes())

	opts.SetWALTtlSeconds(52)
	require.EqualValues(t, 52, opts.GetWALTtlSeconds())

	opts.SetWalSizeLimitMb(540)
	require.EqualValues(t, 540, opts.GetWalSizeLimitMb())

	opts.SetManifestPreallocationSize(5 << 10)
	require.EqualValues(t, 5<<10, opts.GetManifestPreallocationSize())

	opts.SetAllowMmapReads(true)
	require.EqualValues(t, true, opts.AllowMmapReads())

	require.EqualValues(t, false, opts.AllowMmapWrites())
	opts.SetAllowMmapWrites(true)
	require.EqualValues(t, true, opts.AllowMmapWrites())

	opts.SetUseDirectReads(true)
	require.EqualValues(t, true, opts.UseDirectReads())

	opts.SetUseDirectIOForFlushAndCompaction(true)
	require.EqualValues(t, true, opts.UseDirectIOForFlushAndCompaction())

	opts.SetIsFdCloseOnExec(true)
	require.EqualValues(t, true, opts.IsFdCloseOnExec())

	opts.SetSkipLogErrorOnRecovery(true)
	require.EqualValues(t, true, opts.SkipLogErrorOnRecovery())

	opts.SetStatsDumpPeriodSec(79)
	require.EqualValues(t, 79, opts.GetStatsDumpPeriodSec())

	opts.SetStatsPersistPeriodSec(97)
	require.EqualValues(t, 97, opts.GetStatsPersistPeriodSec())

	opts.SetAdviseRandomOnOpen(true)
	require.EqualValues(t, true, opts.AdviseRandomOnOpen())

	opts.SetAccessHintOnCompactionStart(SequentialCompactionAccessPattern)
	require.EqualValues(t, SequentialCompactionAccessPattern, opts.GetAccessHintOnCompactionStart())

	opts.SetDbWriteBufferSize(1 << 30)
	require.EqualValues(t, 1<<30, opts.GetDbWriteBufferSize())

	opts.SetUseAdaptiveMutex(true)
	require.EqualValues(t, true, opts.UseAdaptiveMutex())

	opts.SetBytesPerSync(68 << 10)
	require.EqualValues(t, 68<<10, opts.GetBytesPerSync())

	opts.SetWALBytesPerSync(69 << 10)
	require.EqualValues(t, 69<<10, opts.GetWALBytesPerSync())

	opts.SetWritableFileMaxBufferSize(9 << 20)
	require.EqualValues(t, 9<<20, opts.GetWritableFileMaxBufferSize())

	opts.SetAllowConcurrentMemtableWrites(true)
	require.EqualValues(t, true, opts.AllowConcurrentMemtableWrites())

	opts.SetEnableWriteThreadAdaptiveYield(true)
	require.EqualValues(t, true, opts.EnabledWriteThreadAdaptiveYield())

	opts.SetMaxSequentialSkipInIterations(199)
	require.EqualValues(t, 199, opts.GetMaxSequentialSkipInIterations())

	opts.SetDisableAutoCompactions(true)
	require.EqualValues(t, true, opts.DisabledAutoCompactions())

	opts.SetOptimizeFiltersForHits(true)
	require.EqualValues(t, true, opts.OptimizeFiltersForHits())

	opts.SetDeleteObsoleteFilesPeriodMicros(1234)
	require.EqualValues(t, 1234, opts.GetDeleteObsoleteFilesPeriodMicros())

	opts.SetMemTablePrefixBloomSizeRatio(0.3)
	require.EqualValues(t, 0.3, opts.GetMemTablePrefixBloomSizeRatio())

	opts.SetMaxCompactionBytes(111222)
	require.EqualValues(t, 111222, opts.GetMaxCompactionBytes())

	// cloning
	cl := opts.Clone()
	require.EqualValues(t, 5, cl.GetTableCacheNumshardbits())
}
