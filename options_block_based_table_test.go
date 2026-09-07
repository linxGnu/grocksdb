package grocksdb

import (
	"testing"
)

func TestBBT(t *testing.T) {
	t.Parallel()

	b := NewDefaultBlockBasedTableOptions()
	defer b.Destroy()

	b.SetBlockSize(123)
	b.SetOptimizeFiltersForMemory(true)
	b.SetTopLevelIndexPinningTier(KFallbackPinningTier)
	b.SetPartitionPinningTier(KNonePinningTier)
	b.SetUnpartitionedPinningTier(KAllPinningTier)
	b.SetFormatVersion(4)
	b.SetSeparateKeyValueInDataBlock(true)
	b.SetDataBlockIndexType(KDataBlockIndexTypeBinarySearch)
	b.SetIndexBlockSearchType(KIndexBlockSearchTypeBinary)
	b.SetBlockAlign(true)
}
