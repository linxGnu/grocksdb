package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSTFileManager(t *testing.T) {
	t.Parallel()

	env := NewDefaultEnv()
	defer env.Destroy()

	s := NewSSTFileManager(env)
	defer s.Destroy()

	s.SetMaxAllowedSpaceUsage(100)
	s.SetCompactionBufferSize(1000)

	require.False(t, s.IsMaxAllowedSpaceReached())
	require.False(t, s.IsMaxAllowedSpaceReachedIncludingCompactions())
	require.EqualValues(t, 0, s.GetTotalSize())

	s.SetDeleteRateBytesPerSecond(123)
	require.EqualValues(t, 123, s.GetDeleteRateBytesPerSecond())

	s.SetMaxTrashDBRatio(12.1)
	require.EqualValues(t, 12.1, s.GetMaxTrashDBRatio())

	require.EqualValues(t, 0, s.GetTotalTrashSize())
}
