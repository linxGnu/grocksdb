package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterPolicy(t *testing.T) {
	var (
		givenKeys          = [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
		givenFilter        = []byte("key")
		createFilterCalled = false
		keyMayMatchCalled  = false
	)
	policy := &mockFilterPolicy{
		createFilter: func(keys [][]byte) []byte {
			createFilterCalled = true
			require.EqualValues(t, keys, givenKeys)
			return givenFilter
		},
		keyMayMatch: func(key, filter []byte) bool {
			keyMayMatchCalled = true
			require.EqualValues(t, key, givenKeys[0])
			require.EqualValues(t, filter, givenFilter)
			return true
		},
	}

	db := newTestDB(t, "TestFilterPolicy", func(opts *Options) {
		blockOpts := NewDefaultBlockBasedTableOptions()
		blockOpts.SetFilterPolicy(policy)
		opts.SetBlockBasedTableFactory(blockOpts)
	})
	defer db.Close()

	// insert keys
	wo := NewDefaultWriteOptions()
	for _, k := range givenKeys {
		require.Nil(t, db.Put(wo, k, []byte("val")))
	}

	// flush to trigger the filter creation
	require.Nil(t, db.Flush(NewDefaultFlushOptions()))
	require.True(t, createFilterCalled)

	// test key may match call
	ro := NewDefaultReadOptions()
	v1, err := db.Get(ro, givenKeys[0])
	defer v1.Free()
	require.Nil(t, err)
	require.True(t, keyMayMatchCalled)
}

type mockFilterPolicy struct {
	createFilter func(keys [][]byte) []byte
	keyMayMatch  func(key, filter []byte) bool
}

func (m *mockFilterPolicy) Name() string { return "grocksdb.test" }

func (m *mockFilterPolicy) CreateFilter(keys [][]byte) []byte {
	return m.createFilter(keys)
}

func (m *mockFilterPolicy) KeyMayMatch(key, filter []byte) bool {
	return m.keyMayMatch(key, filter)
}

func (m *mockFilterPolicy) Destroy() {}
