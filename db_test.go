package grocksdb

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestOpenDb(t *testing.T) {
	db := newTestDB(t, "TestOpenDb", nil)
	defer db.Close()
	ensure.DeepEqual(t, "0", db.GetProperty("rocksdb.num-immutable-mem-table"))
	v, success := db.GetIntProperty("rocksdb.num-immutable-mem-table")
	ensure.DeepEqual(t, uint64(0), v)
	ensure.True(t, success)
}

func TestDBCRUD(t *testing.T) {
	db := newTestDB(t, "TestDBGet", nil)
	defer db.Close()

	var (
		givenKey  = []byte("hello")
		givenVal1 = []byte("")
		givenVal2 = []byte("world1")
		wo        = NewDefaultWriteOptions()
		ro        = NewDefaultReadOptions()
	)

	// create
	ensure.Nil(t, db.Put(wo, givenKey, givenVal1))

	// retrieve
	v1, err := db.Get(ro, givenKey)
	defer v1.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, v1.Data(), givenVal1)

	// retrieve bytes
	_v1, err := db.GetBytes(ro, givenKey)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, _v1, givenVal1)

	// update
	ensure.Nil(t, db.Put(wo, givenKey, givenVal2))
	v2, err := db.Get(ro, givenKey)
	defer v2.Free()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, v2.Data(), givenVal2)

	// retrieve pinned
	v3, err := db.GetPinned(ro, givenKey)
	defer v3.Destroy()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, v3.Data(), givenVal2)

	// delete
	ensure.Nil(t, db.Delete(wo, givenKey))
	v4, err := db.Get(ro, givenKey)
	ensure.Nil(t, err)
	ensure.True(t, v4.Data() == nil)

	// retrieve missing pinned
	v5, err := db.GetPinned(ro, givenKey)
	defer v5.Destroy()
	ensure.Nil(t, err)
	ensure.True(t, v5.Data() == nil)
}

func TestDBCRUDDBPaths(t *testing.T) {
	names := make([]string, 4)
	targetSizes := make([]uint64, len(names))

	for i := range names {
		names[i] = "TestDBGet_" + strconv.FormatInt(int64(i), 10)
		targetSizes[i] = uint64(1024 * 1024 * (i + 1))
	}

	db := newTestDBPathNames(t, "TestDBGet", names, targetSizes, nil)
	defer db.Close()

	var (
		givenKey  = []byte("hello")
		givenVal1 = []byte("")
		givenVal2 = []byte("world1")
		givenVal3 = []byte("world2")
		wo        = NewDefaultWriteOptions()
		ro        = NewDefaultReadOptions()
	)

	// retrieve before create
	noexist, err := db.Get(ro, givenKey)
	defer noexist.Free()
	ensure.Nil(t, err)
	ensure.False(t, noexist.Exists())
	ensure.DeepEqual(t, noexist.Data(), []byte(nil))

	// create
	ensure.Nil(t, db.Put(wo, givenKey, givenVal1))

	// retrieve
	v1, err := db.Get(ro, givenKey)
	defer v1.Free()
	ensure.Nil(t, err)
	ensure.True(t, v1.Exists())
	ensure.DeepEqual(t, v1.Data(), givenVal1)

	// update
	ensure.Nil(t, db.Put(wo, givenKey, givenVal2))
	v2, err := db.Get(ro, givenKey)
	defer v2.Free()
	ensure.Nil(t, err)
	ensure.True(t, v2.Exists())
	ensure.DeepEqual(t, v2.Data(), givenVal2)

	// update
	ensure.Nil(t, db.Put(wo, givenKey, givenVal3))
	v3, err := db.Get(ro, givenKey)
	defer v3.Free()
	ensure.Nil(t, err)
	ensure.True(t, v3.Exists())
	ensure.DeepEqual(t, v3.Data(), givenVal3)

	// delete
	ensure.Nil(t, db.Delete(wo, givenKey))
	v4, err := db.Get(ro, givenKey)
	defer v4.Free()
	ensure.Nil(t, err)
	ensure.False(t, v4.Exists())
	ensure.DeepEqual(t, v4.Data(), []byte(nil))
}

func newTestDB(t *testing.T, name string, applyOpts func(opts *Options)) *DB {
	dir, err := ioutil.TempDir("", "gorocksdb-"+name)
	ensure.Nil(t, err)

	opts := NewDefaultOptions()
	// test the ratelimiter
	rateLimiter := NewRateLimiter(1024, 100*1000, 10)
	opts.SetRateLimiter(rateLimiter)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(ZSTDCompression)
	if applyOpts != nil {
		applyOpts(opts)
	}
	db, err := OpenDb(opts, dir)
	ensure.Nil(t, err)

	return db
}

func newTestDBMultiCF(t *testing.T, name string, columns []string, applyOpts func(opts *Options)) (db *DB, cfh []*ColumnFamilyHandle, cleanup func()) {
	dir, err := ioutil.TempDir("", "gorocksdb-"+name)
	ensure.Nil(t, err)

	opts := NewDefaultOptions()
	opts.SetCreateIfMissingColumnFamilies(true)
	opts.SetCreateIfMissing(true)
	opts.SetCompression(ZSTDCompression)
	options := make([]*Options, len(columns))
	for i := range options {
		options[i] = opts
	}

	db, cfh, err = OpenDbColumnFamilies(opts, dir, columns, options)
	ensure.Nil(t, err)
	cleanup = func() {
		for _, cf := range cfh {
			cf.Destroy()
		}
		db.Close()
	}
	return db, cfh, cleanup
}

func newTestDBPathNames(t *testing.T, name string, names []string, targetSizes []uint64, applyOpts func(opts *Options)) *DB {
	ensure.DeepEqual(t, len(targetSizes), len(names))
	ensure.NotDeepEqual(t, len(names), 0)

	dir, err := ioutil.TempDir("", "gorocksdb-"+name)
	ensure.Nil(t, err)

	paths := make([]string, len(names))
	for i, name := range names {
		dir, err := ioutil.TempDir("", "gorocksdb-"+name)
		ensure.Nil(t, err)
		paths[i] = dir
	}

	dbpaths := NewDBPathsFromData(paths, targetSizes)
	defer DestroyDBPaths(dbpaths)

	opts := NewDefaultOptions()
	opts.SetDBPaths(dbpaths)
	// test the ratelimiter
	rateLimiter := NewRateLimiter(1024, 100*1000, 10)
	opts.SetRateLimiter(rateLimiter)
	opts.SetCreateIfMissing(true)
	if applyOpts != nil {
		applyOpts(opts)
	}
	db, err := OpenDb(opts, dir)
	ensure.Nil(t, err)

	return db
}

func TestDBMultiGet(t *testing.T) {
	db := newTestDB(t, "TestDBMultiGet", nil)
	defer db.Close()

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
	ensure.Nil(t, db.Put(wo, givenKey1, givenVal1))
	ensure.Nil(t, db.Put(wo, givenKey2, givenVal2))
	ensure.Nil(t, db.Put(wo, givenKey3, givenVal3))

	// retrieve
	values, err := db.MultiGet(ro, []byte("noexist"), givenKey1, givenKey2, givenKey3)
	defer values.Destroy()
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(values), 4)

	ensure.DeepEqual(t, values[0].Data(), []byte(nil))
	ensure.DeepEqual(t, values[1].Data(), givenVal1)
	ensure.DeepEqual(t, values[2].Data(), givenVal2)
	ensure.DeepEqual(t, values[3].Data(), givenVal3)
}

func TestDBGetApproximateSizes(t *testing.T) {
	db := newTestDB(t, "TestDBGetApproximateSizes", nil)
	defer db.Close()

	// no ranges
	sizes := db.GetApproximateSizes(nil)
	ensure.DeepEqual(t, len(sizes), 0)

	// range will nil start and limit
	sizes = db.GetApproximateSizes([]Range{{Start: nil, Limit: nil}})
	ensure.DeepEqual(t, sizes, []uint64{0})

	// valid range
	sizes = db.GetApproximateSizes([]Range{{Start: []byte{0x00}, Limit: []byte{0xFF}}})
	ensure.DeepEqual(t, sizes, []uint64{0})
}

func TestDBGetApproximateSizesCF(t *testing.T) {
	db := newTestDB(t, "TestDBGetApproximateSizesCF", nil)
	defer db.Close()

	o := NewDefaultOptions()

	cf, err := db.CreateColumnFamily(o, "other")
	ensure.Nil(t, err)

	// no ranges
	sizes := db.GetApproximateSizesCF(cf, nil)
	ensure.DeepEqual(t, len(sizes), 0)

	// range will nil start and limit
	sizes = db.GetApproximateSizesCF(cf, []Range{{Start: nil, Limit: nil}})
	ensure.DeepEqual(t, sizes, []uint64{0})

	// valid range
	sizes = db.GetApproximateSizesCF(cf, []Range{{Start: []byte{0x00}, Limit: []byte{0xFF}}})
	ensure.DeepEqual(t, sizes, []uint64{0})
}
