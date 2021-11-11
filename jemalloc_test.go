// +build !testing,jemalloc

package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemAlloc(t *testing.T) {
	m, err := CreateJemallocNodumpAllocator()
	require.NoError(t, err)
	m.Destroy()
}
