package grocksdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteBatch(t *testing.T) {
	db := newTestDB(t, "TestWriteBatch", nil)
	defer db.Close()

	var (
		givenKey1 = []byte("key1")
		givenVal1 = []byte("val1")
		givenKey2 = []byte("key2")
	)
	wo := NewDefaultWriteOptions()
	require.Nil(t, db.Put(wo, givenKey2, []byte("foo")))

	// create and fill the write batch
	wb := NewWriteBatch()
	defer wb.Destroy()
	wb.Put(givenKey1, givenVal1)
	wb.Delete(givenKey2)
	require.EqualValues(t, wb.Count(), 2)

	// perform the batch
	require.Nil(t, db.Write(wo, wb))

	// check changes
	ro := NewDefaultReadOptions()
	v1, err := db.Get(ro, givenKey1)
	defer v1.Free()
	require.Nil(t, err)
	require.EqualValues(t, v1.Data(), givenVal1)

	v2, err := db.Get(ro, givenKey2)
	defer v2.Free()
	require.Nil(t, err)
	require.True(t, v2.Data() == nil)

	// DeleteRange test
	wb.Clear()
	wb.DeleteRange(givenKey1, givenKey2)

	// perform the batch
	require.Nil(t, db.Write(wo, wb))

	v1, err = db.Get(ro, givenKey1)
	defer v1.Free()
	require.Nil(t, err)
	require.True(t, v1.Data() == nil)
}

func TestWriteBatchIterator(t *testing.T) {
	db := newTestDB(t, "TestWriteBatchIterator", nil)
	defer db.Close()

	var (
		givenKey1 = []byte("key1")
		givenVal1 = []byte("val1")
		givenKey2 = []byte("key2")
	)
	// create and fill the write batch
	wb := NewWriteBatch()
	defer wb.Destroy()
	wb.Put(givenKey1, givenVal1)
	wb.Delete(givenKey2)
	require.EqualValues(t, wb.Count(), 2)

	// iterate over the batch
	iter := wb.NewIterator()
	require.True(t, iter.Next())
	record := iter.Record()
	require.EqualValues(t, record.Type, WriteBatchValueRecord)
	require.EqualValues(t, record.Key, givenKey1)
	require.EqualValues(t, record.Value, givenVal1)

	require.True(t, iter.Next())
	record = iter.Record()
	require.EqualValues(t, record.Type, WriteBatchDeletionRecord)
	require.EqualValues(t, record.Key, givenKey2)

	// there shouldn't be any left
	require.False(t, iter.Next())
}
