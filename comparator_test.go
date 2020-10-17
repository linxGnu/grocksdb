package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComparator(t *testing.T) {
	db := newTestDB(t, "TestComparator", func(opts *Options) {
		opts.SetComparator(&testBytesReverseComparator{})
	})
	defer db.Close()

	// insert keys
	givenKeys := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
	wo := NewDefaultWriteOptions()
	for _, k := range givenKeys {
		require.Nil(t, db.Put(wo, k, []byte("val")))
	}

	// create a iterator to collect the keys
	ro := NewDefaultReadOptions()
	iter := db.NewIterator(ro)
	defer iter.Close()

	// we seek to the last key and iterate in reverse order
	// to match given keys
	var actualKeys [][]byte
	for iter.SeekToLast(); iter.Valid(); iter.Prev() {
		key := make([]byte, 4)
		copy(key, iter.Key().Data())
		actualKeys = append(actualKeys, key)
	}
	require.Nil(t, iter.Err())

	// ensure that the order is correct
	require.EqualValues(t, actualKeys, givenKeys)
}
