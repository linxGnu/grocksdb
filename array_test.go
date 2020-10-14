package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToCSlice(t *testing.T) {
	v, err := byteSlicesToCSlices(nil)
	require.Nil(t, v)
	require.Nil(t, err)
}
