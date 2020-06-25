package grocksdb

import (
	"testing"
)

func TestOptions(t *testing.T) {
	opts := NewDefaultOptions()
	opts.SetDumpMallocStats(true)
	opts.SetMemtableWholeKeyFiltering(true)
	opts.Destroy()
}
