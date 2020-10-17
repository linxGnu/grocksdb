package grocksdb

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumnFamilyOpen(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyOpen")
	require.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(LZ4Compression)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	require.Nil(t, err)
	defer db.Close()
	require.EqualValues(t, len(cfh), 2)
	cfh[0].Destroy()
	cfh[1].Destroy()

	actualNames, err := ListColumnFamilies(opts, dir)
	require.Nil(t, err)
	require.EqualValues(t, actualNames, givenNames)
}

func TestColumnFamilyCreateDrop(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyCreate")
	require.Nil(t, err)

	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(LZ4HCCompression)
	db, err := OpenDb(opts, dir)
	require.Nil(t, err)
	defer db.Close()
	cf, err := db.CreateColumnFamily(opts, "guide")
	require.Nil(t, err)
	defer cf.Destroy()

	actualNames, err := ListColumnFamilies(opts, dir)
	require.Nil(t, err)
	require.EqualValues(t, actualNames, []string{"default", "guide"})

	require.Nil(t, db.DropColumnFamily(cf))

	actualNames, err = ListColumnFamilies(opts, dir)
	require.Nil(t, err)
	require.EqualValues(t, actualNames, []string{"default"})
}

func TestColumnFamilyBatchPutGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	require.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	require.Nil(t, err)
	defer db.Close()
	require.EqualValues(t, len(cfh), 2)
	defer cfh[0].Destroy()
	defer cfh[1].Destroy()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	givenKey0 := []byte("hello0")
	givenVal0 := []byte("world0")
	givenKey1 := []byte("hello1")
	givenVal1 := []byte("world1")

	b0 := NewWriteBatch()
	defer b0.Destroy()
	b0.PutCF(cfh[0], givenKey0, givenVal0)
	require.Nil(t, db.Write(wo, b0))
	actualVal0, err := db.GetCF(ro, cfh[0], givenKey0)
	defer actualVal0.Free()
	require.Nil(t, err)
	require.EqualValues(t, actualVal0.Data(), givenVal0)

	b1 := NewWriteBatch()
	defer b1.Destroy()
	b1.PutCF(cfh[1], givenKey1, givenVal1)
	require.Nil(t, db.Write(wo, b1))
	actualVal1, err := db.GetCF(ro, cfh[1], givenKey1)
	defer actualVal1.Free()
	require.Nil(t, err)
	require.EqualValues(t, actualVal1.Data(), givenVal1)

	actualVal, err := db.GetCF(ro, cfh[0], givenKey1)
	require.Nil(t, err)
	require.EqualValues(t, actualVal.Size(), 0)
	actualVal, err = db.GetCF(ro, cfh[1], givenKey0)
	require.Nil(t, err)
	require.EqualValues(t, actualVal.Size(), 0)

	{
		v := db.KeyMayExistsCF(ro, cfh[0], givenKey0, "")
		defer v.Free()
		require.True(t, v.Size() > 0)
	}

	// trigger flush
	require.Nil(t, db.FlushCF(cfh[0], NewDefaultFlushOptions()))
}

func TestColumnFamilyPutGetDelete(t *testing.T) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	require.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(SnappyCompression)
	db, cfh, err := OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	require.Nil(t, err)
	defer db.Close()
	require.EqualValues(t, len(cfh), 2)
	defer cfh[0].Destroy()
	defer cfh[1].Destroy()

	wo := NewDefaultWriteOptions()
	defer wo.Destroy()
	ro := NewDefaultReadOptions()
	defer ro.Destroy()

	givenKey0 := []byte("hello0")
	givenVal0 := []byte("world0")
	givenKey1 := []byte("hello1")
	givenVal1 := []byte("world1")

	{
		require.Nil(t, db.PutCF(wo, cfh[0], givenKey0, givenVal0))
		actualVal0, err := db.GetCF(ro, cfh[0], givenKey0)
		defer actualVal0.Free()
		require.Nil(t, err)
		require.EqualValues(t, actualVal0.Data(), givenVal0)

		require.Nil(t, db.PutCF(wo, cfh[1], givenKey1, givenVal1))
		actualVal1, err := db.GetCF(ro, cfh[1], givenKey1)
		defer actualVal1.Free()
		require.Nil(t, err)
		require.EqualValues(t, actualVal1.Data(), givenVal1)

		actualVal, err := db.GetCF(ro, cfh[0], givenKey1)
		require.Nil(t, err)
		require.EqualValues(t, actualVal.Size(), 0)
		actualVal, err = db.GetCF(ro, cfh[1], givenKey0)
		require.Nil(t, err)
		require.EqualValues(t, actualVal.Size(), 0)

		require.Nil(t, db.DeleteCF(wo, cfh[0], givenKey0))
		actualVal, err = db.GetCF(ro, cfh[0], givenKey0)
		require.Nil(t, err)
		require.EqualValues(t, actualVal.Size(), 0)

		{
			v := db.KeyMayExistsCF(ro, cfh[0], givenKey0, "")
			defer v.Free()
		}
	}

	{
		require.Nil(t, db.PutCF(wo, cfh[0], givenKey0, givenVal0))
		actualVal0, err := db.GetCF(ro, cfh[0], givenKey0)
		defer actualVal0.Free()
		require.Nil(t, err)
		require.EqualValues(t, actualVal0.Data(), givenVal0)

		require.Nil(t, db.DeleteRangeCF(wo, cfh[0], givenKey0, givenKey1))
		actualVal, err := db.GetCF(ro, cfh[0], givenKey0)
		require.Nil(t, err)
		require.EqualValues(t, actualVal.Size(), 0)

		actualVal1, err := db.GetCF(ro, cfh[1], givenKey1)
		defer actualVal1.Free()
		require.Nil(t, err)
		require.EqualValues(t, actualVal1.Data(), givenVal1)
	}
}

func newTestDBCF(t *testing.T, name string) (db *DB, cfh []*ColumnFamilyHandle, cleanup func()) {
	dir, err := ioutil.TempDir("", "gorocksdb-TestColumnFamilyPutGet")
	require.Nil(t, err)

	givenNames := []string{"default", "guide"}
	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(ZLibCompression)
	db, cfh, err = OpenDbColumnFamilies(opts, dir, givenNames, []*Options{opts, opts})
	require.Nil(t, err)
	cleanup = func() {
		for _, cf := range cfh {
			cf.Destroy()
		}
		db.Close()
	}
	return db, cfh, cleanup
}

func TestColumnFamilyMultiGet(t *testing.T) {
	db, cfh, cleanup := newTestDBCF(t, "TestDBMultiGet")
	defer cleanup()

	var (
		givenKey1 = []byte("hello1")
		givenKey2 = []byte("hello2")
		givenKey3 = []byte("hello3")
		givenVal1 = []byte("world1")
		givenVal2 = []byte("world2")
		givenVal3 = []byte("world3")
		wo        = NewDefaultWriteOptions()
		ro        = NewDefaultReadOptions()
	)

	// create
	require.Nil(t, db.PutCF(wo, cfh[0], givenKey1, givenVal1))
	require.Nil(t, db.PutCF(wo, cfh[1], givenKey2, givenVal2))
	require.Nil(t, db.PutCF(wo, cfh[1], givenKey3, givenVal3))

	// column family 0 only has givenKey1
	values, err := db.MultiGetCF(ro, cfh[0], []byte("noexist"), givenKey1, givenKey2, givenKey3)
	defer values.Destroy()
	require.Nil(t, err)
	require.EqualValues(t, len(values), 4)

	require.EqualValues(t, values[0].Data(), []byte(nil))
	require.EqualValues(t, values[1].Data(), givenVal1)
	require.EqualValues(t, values[2].Data(), []byte(nil))
	require.EqualValues(t, values[3].Data(), []byte(nil))

	// column family 1 only has givenKey2 and givenKey3
	values, err = db.MultiGetCF(ro, cfh[1], []byte("noexist"), givenKey1, givenKey2, givenKey3)
	defer values.Destroy()
	require.Nil(t, err)
	require.EqualValues(t, len(values), 4)

	require.EqualValues(t, values[0].Data(), []byte(nil))
	require.EqualValues(t, values[1].Data(), []byte(nil))
	require.EqualValues(t, values[2].Data(), givenVal2)
	require.EqualValues(t, values[3].Data(), givenVal3)

	// getting them all from the right CF should return them all
	values, err = db.MultiGetCFMultiCF(ro,
		ColumnFamilyHandles{cfh[0], cfh[1], cfh[1]},
		[][]byte{givenKey1, givenKey2, givenKey3},
	)
	defer values.Destroy()
	require.Nil(t, err)
	require.EqualValues(t, len(values), 3)

	require.EqualValues(t, values[0].Data(), givenVal1)
	require.EqualValues(t, values[1].Data(), givenVal2)
	require.EqualValues(t, values[2].Data(), givenVal3)
}
