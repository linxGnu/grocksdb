package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemAlloc(t *testing.T) {
	_, err := CreateJemallocNodumpAllocator()
	require.Error(t, err)
}
