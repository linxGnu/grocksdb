package grocksdb

import (
	"testing"
)

func TestBBT(t *testing.T) {
	b := NewDefaultBlockBasedTableOptions()
	defer b.Destroy()

	b.SetBlockSize(123)
	b.SetOptimizeFiltersForMemory(true)
}
