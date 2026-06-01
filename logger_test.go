package grocksdb

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallbackLogger(t *testing.T) {
	t.Parallel()

	var (
		mu      sync.Mutex
		count   int64
		lastMsg string
		lastLvl InfoLogLevel
	)

	logger := NewCallbackLogger(DebugInfoLogLevel, func(level InfoLogLevel, msg string) {
		atomic.AddInt64(&count, 1)
		mu.Lock()
		lastMsg = msg
		lastLvl = level
		mu.Unlock()
	})

	db := newTestDB(t, func(opts *Options) {
		opts.SetInfoLog(logger)
	})
	defer db.Close()

	// Generate some log activity.
	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	require.NoError(t, db.Put(wo, []byte("foo"), []byte("bar")))
	fo := NewDefaultFlushOptions()
	defer fo.Destroy()
	require.NoError(t, db.Flush(fo))

	require.Greater(t, atomic.LoadInt64(&count), int64(0))

	mu.Lock()
	require.NotEmpty(t, lastMsg)
	_ = lastLvl
	mu.Unlock()
}
