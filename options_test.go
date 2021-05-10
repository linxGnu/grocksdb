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

	require.EqualValues(t, PointInTimeRecovery, opts.GetWALRecoveryMode())
	opts.SetWALRecoveryMode(SkipAnyCorruptedRecordsRecovery)
	require.EqualValues(t, SkipAnyCorruptedRecordsRecovery, opts.GetWALRecoveryMode())

	require.EqualValues(t, 2, opts.GetMaxBackgroundJobs())
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

	require.EqualValues(t, uint64(0x40000000), opts.GetMaxManifestFileSize())
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

	require.EqualValues(t, 4<<20, opts.GetManifestPreallocationSize())
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

	opts.SetMemtableHugePageSize(223344)
	require.EqualValues(t, 223344, opts.GetMemtableHugePageSize())

	opts.SetMaxSuccessiveMerges(99)
	require.EqualValues(t, 99, opts.GetMaxSuccessiveMerges())

	opts.SetBloomLocality(5)
	require.EqualValues(t, 5, opts.GetBloomLocality())

	require.EqualValues(t, false, opts.InplaceUpdateSupport())
	opts.SetInplaceUpdateSupport(true)
	require.EqualValues(t, true, opts.InplaceUpdateSupport())

	require.EqualValues(t, 10000, opts.GetInplaceUpdateNumLocks())
	opts.SetInplaceUpdateNumLocks(8)
	require.EqualValues(t, 8, opts.GetInplaceUpdateNumLocks())

	opts.SetReportBackgroundIOStats(true)
	require.EqualValues(t, true, opts.ReportBackgroundIOStats())

	opts.SetMaxTotalWalSize(10 << 30)
	require.EqualValues(t, 10<<30, opts.GetMaxTotalWalSize())

	opts.SetBottommostCompression(ZLibCompression)
	require.EqualValues(t, ZLibCompression, opts.GetBottommostCompression())

	require.EqualValues(t, SnappyCompression, opts.GetCompression())
	opts.SetCompression(LZ4Compression)
	require.EqualValues(t, LZ4Compression, opts.GetCompression())

	require.EqualValues(t, LevelCompactionStyle, opts.GetCompactionStyle())
	opts.SetCompactionStyle(UniversalCompactionStyle)
	require.EqualValues(t, UniversalCompactionStyle, opts.GetCompactionStyle())

	require.EqualValues(t, false, opts.IsAtomicFlush())
	opts.SetAtomicFlush(true)
	require.EqualValues(t, true, opts.IsAtomicFlush())

	require.EqualValues(t, false, opts.CreateIfMissing())
	opts.SetCreateIfMissing(true)
	require.EqualValues(t, true, opts.CreateIfMissing())

	require.EqualValues(t, false, opts.CreateIfMissingColumnFamilies())
	opts.SetCreateIfMissingColumnFamilies(true)
	require.EqualValues(t, true, opts.CreateIfMissingColumnFamilies())

	opts.SetErrorIfExists(true)
	require.EqualValues(t, true, opts.ErrorIfExists())

	opts.SetParanoidChecks(true)
	require.EqualValues(t, true, opts.ParanoidChecks())

	require.EqualValues(t, InfoInfoLogLevel, opts.GetInfoLogLevel())
	opts.SetInfoLogLevel(WarnInfoLogLevel)
	require.EqualValues(t, WarnInfoLogLevel, opts.GetInfoLogLevel())

	require.EqualValues(t, 64<<20, opts.GetWriteBufferSize())
	opts.SetWriteBufferSize(1 << 19)
	require.EqualValues(t, 1<<19, opts.GetWriteBufferSize())

	require.EqualValues(t, 2, opts.GetMaxWriteBufferNumber())
	opts.SetMaxWriteBufferNumber(15)
	require.EqualValues(t, 15, opts.GetMaxWriteBufferNumber())

	require.EqualValues(t, 1, opts.GetMinWriteBufferNumberToMerge())
	opts.SetMinWriteBufferNumberToMerge(2)
	require.EqualValues(t, 2, opts.GetMinWriteBufferNumberToMerge())

	require.EqualValues(t, -1, opts.GetMaxOpenFiles())
	opts.SetMaxOpenFiles(999)
	require.EqualValues(t, 999, opts.GetMaxOpenFiles())

	require.EqualValues(t, 16, opts.GetMaxFileOpeningThreads())
	opts.SetMaxFileOpeningThreads(21)
	require.EqualValues(t, 21, opts.GetMaxFileOpeningThreads())

	opts.SetCompressionPerLevel([]CompressionType{ZLibCompression, SnappyCompression})

	opts.SetEnv(NewMemEnv())
	opts.SetEnv(NewDefaultEnv())

	opts.IncreaseParallelism(8)

	opts.OptimizeForPointLookup(19 << 20)

	opts.OptimizeLevelStyleCompaction(10 << 20)

	opts.OptimizeUniversalStyleCompaction(20 << 20)

	require.EqualValues(t, true, opts.AllowConcurrentMemtableWrites())
	opts.SetAllowConcurrentMemtableWrites(false)
	require.EqualValues(t, false, opts.AllowConcurrentMemtableWrites())

	opts.SetCompressionOptionsZstdMaxTrainBytes(123 << 20)
	opts.SetCompressionOptionsMaxDictBufferBytes(213 << 10)
	opts.SetBottommostCompressionOptionsZstdMaxTrainBytes(234<<20, true)
	opts.SetBottommostCompressionOptionsMaxDictBufferBytes(312<<10, true)

	opts.SetBottommostCompressionOptions(NewDefaultCompressionOptions(), true)
	opts.SetCompressionOptions(NewDefaultCompressionOptions())
	opts.SetMinLevelToCompress(2)

	require.EqualValues(t, 7, opts.GetNumLevels())
	opts.SetNumLevels(8)
	require.EqualValues(t, 8, opts.GetNumLevels())

	require.EqualValues(t, 2, opts.GetLevel0FileNumCompactionTrigger())
	opts.SetLevel0FileNumCompactionTrigger(14)
	require.EqualValues(t, 14, opts.GetLevel0FileNumCompactionTrigger())

	require.EqualValues(t, 20, opts.GetLevel0SlowdownWritesTrigger())
	opts.SetLevel0SlowdownWritesTrigger(17)
	require.EqualValues(t, 17, opts.GetLevel0SlowdownWritesTrigger())

	require.EqualValues(t, 36, opts.GetLevel0StopWritesTrigger())
	opts.SetLevel0StopWritesTrigger(47)
	require.EqualValues(t, 47, opts.GetLevel0StopWritesTrigger())

	opts.SetMaxMemCompactionLevel(4)

	require.EqualValues(t, uint64(0x140000), opts.GetTargetFileSizeBase())
	opts.SetTargetFileSizeBase(41 << 20)
	require.EqualValues(t, 41<<20, opts.GetTargetFileSizeBase())

	require.EqualValues(t, 1, opts.GetTargetFileSizeMultiplier())
	opts.SetTargetFileSizeMultiplier(3)
	require.EqualValues(t, 3, opts.GetTargetFileSizeMultiplier())

	require.EqualValues(t, 10<<20, opts.GetMaxBytesForLevelBase())
	opts.SetMaxBytesForLevelBase(1 << 30)
	require.EqualValues(t, 1<<30, opts.GetMaxBytesForLevelBase())

	require.EqualValues(t, 10, opts.GetMaxBytesForLevelMultiplier())
	opts.SetMaxBytesForLevelMultiplier(12)
	require.EqualValues(t, 12, opts.GetMaxBytesForLevelMultiplier())

	require.EqualValues(t, 1, opts.GetMaxSubcompactions())
	opts.SetMaxSubcompactions(3)
	require.EqualValues(t, 3, opts.GetMaxSubcompactions())

	opts.SetMaxBytesForLevelMultiplierAdditional([]int{2 << 20})

	opts.SetPurgeRedundantKvsWhileFlush(true)

	opts.SetDbLogDir("./abc")
	opts.SetWalDir("../asdf")

	require.EqualValues(t, false, opts.EnabledPipelinedWrite())
	opts.SetEnablePipelinedWrite(true)
	require.EqualValues(t, true, opts.EnabledPipelinedWrite())

	require.EqualValues(t, false, opts.UnorderedWrite())
	opts.SetUnorderedWrite(true)
	require.EqualValues(t, true, opts.UnorderedWrite())

	opts.EnableStatistics()
	opts.PrepareForBulkLoad()
	opts.SetMemtableVectorRep()
	opts.SetHashLinkListRep(12)
	opts.SetHashSkipListRep(1, 2, 3)
	opts.SetPlainTableFactory(1, 2, 3.1, 12)
	opts.SetUint64AddMergeOperator()
	opts.SetDumpMallocStats(true)
	opts.SetMemtableWholeKeyFiltering(true)

	require.EqualValues(t, false, opts.AllowIngestBehind())
	opts.SetAllowIngestBehind(true)
	require.EqualValues(t, true, opts.AllowIngestBehind())

	require.EqualValues(t, false, opts.SkipStatsUpdateOnDBOpen())
	opts.SetSkipStatsUpdateOnDBOpen(true)
	require.EqualValues(t, true, opts.SkipStatsUpdateOnDBOpen())

	require.EqualValues(t, false, opts.SkipCheckingSSTFileSizesOnDBOpen())
	opts.SetSkipCheckingSSTFileSizesOnDBOpen(true)
	require.EqualValues(t, true, opts.SkipCheckingSSTFileSizesOnDBOpen())

	opts.CompactionReadaheadSize(88 << 20)
	require.EqualValues(t, 88<<20, opts.GetCompactionReadaheadSize())

	opts.SetMaxWriteBufferSizeToMaintain(99 << 19)
	require.EqualValues(t, 99<<19, opts.GetMaxWriteBufferSizeToMaintain())

	// set compaction filter
	opts.SetCompactionFilter(NewNativeCompactionFilter(nil))

	// set comparator
	opts.SetComparator(NewNativeComparator(nil))

	// set merge operator
	opts.SetMergeOperator(NewNativeMergeOperator(nil))

	// get option from string
	_, err := GetOptionsFromString(nil, "abc")
	require.Error(t, err)

	// deprecation soon
	opts.SetTableCacheRemoveScanCountLimit(112)

	opts.SetMaxWriteBufferNumberToMaintain(45)
	require.EqualValues(t, 45, opts.GetMaxWriteBufferNumberToMaintain())

	// cloning
	cl := opts.Clone()
	require.EqualValues(t, 5, cl.GetTableCacheNumshardbits())
}

func TestOptions2(t *testing.T) {
	t.Run("SetUniversalCompactionOpts", func(t *testing.T) {
		opts := NewDefaultOptions()
		defer opts.Destroy()

		opts.SetUniversalCompactionOptions(NewDefaultUniversalCompactionOptions())
	})

	t.Run("SetFifoCompactionOpts", func(t *testing.T) {
		opts := NewDefaultOptions()
		defer opts.Destroy()

		opts.SetFIFOCompactionOptions(NewDefaultFIFOCompactionOptions())
	})

	t.Run("StatisticString", func(t *testing.T) {
		opts := NewDefaultOptions()
		defer opts.Destroy()

		require.Empty(t, opts.GetStatisticsString())
	})
}
