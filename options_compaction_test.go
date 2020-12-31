package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptionCompactions(t *testing.T) {
	co := NewCompactRangeOptions()
	defer co.Destroy()

	require.EqualValues(t, true, co.GetExclusiveManualCompaction())
	co.SetExclusiveManualCompaction(false)
	require.EqualValues(t, false, co.GetExclusiveManualCompaction())

	require.EqualValues(t, KIfHaveCompactionFilter, co.BottommostLevelCompaction())
	co.SetBottommostLevelCompaction(KForce)
	require.EqualValues(t, KForce, co.BottommostLevelCompaction())

	require.EqualValues(t, false, co.ChangeLevel())
	co.SetChangeLevel(true)
	require.EqualValues(t, true, co.ChangeLevel())

	require.EqualValues(t, -1, co.TargetLevel())
	co.SetTargetLevel(2)
	require.EqualValues(t, 2, co.TargetLevel())
}

func TestFifoCompactOption(t *testing.T) {
	fo := NewDefaultFIFOCompactionOptions()
	defer fo.Destroy()

	fo.SetMaxTableFilesSize(2 << 10)
	require.EqualValues(t, 2<<10, fo.GetMaxTableFilesSize())
}

func TestUniversalCompactOption(t *testing.T) {
	uo := NewDefaultUniversalCompactionOptions()
	defer uo.Destroy()

	uo.SetSizeRatio(2)
	require.EqualValues(t, 2, uo.GetSizeRatio())

	uo.SetMinMergeWidth(3)
	require.EqualValues(t, 3, uo.GetMinMergeWidth())

	uo.SetMaxMergeWidth(123)
	require.EqualValues(t, 123, uo.GetMaxMergeWidth())

	uo.SetMaxSizeAmplificationPercent(20)
	require.EqualValues(t, 20, uo.GetMaxSizeAmplificationPercent())

	uo.SetCompressionSizePercent(18)
	require.EqualValues(t, 18, uo.GetCompressionSizePercent())

	uo.SetStopStyle(CompactionStopStyleTotalSize)
	require.EqualValues(t, CompactionStopStyleTotalSize, uo.GetStopStyle())
}
