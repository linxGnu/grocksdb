package grocksdb

import (
	"testing"
)

func TestFilterPolicy(t *testing.T) {
	t.Run("Bloom", func(*testing.T) {
		flt := NewBloomFilter(1.2)
		defer flt.Destroy()
	})

	t.Run("BloomFull", func(*testing.T) {
		flt := NewBloomFilterFull(1.2)
		defer flt.Destroy()
	})

	t.Run("Ribbon", func(*testing.T) {
		flt := NewRibbonFilterPolicy(1.2)
		defer flt.Destroy()
	})

	t.Run("RibbonHybrid", func(*testing.T) {
		flt := NewRibbonHybridFilterPolicy(1.2, 1)
		defer flt.Destroy()
	})
}
