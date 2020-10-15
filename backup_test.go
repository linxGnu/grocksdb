package grocksdb

// import (
// 	"io/ioutil"
// 	"testing"

// 	"github.com/facebookgo/ensure"
// 	"github.com/stretchr/testify/require"
// )

// func TestBackupEngine(t *testing.T) {
// 	db := newTestDB(t, "TestDBBackup", nil)
// 	defer db.Close()

// 	var (
// 		givenKey  = []byte("hello")
// 		givenVal1 = []byte("")
// 		givenVal2 = []byte("world1")
// 		wo        = NewDefaultWriteOptions()
// 		ro        = NewDefaultReadOptions()
// 	)
// 	defer wo.Destroy()
// 	defer ro.Destroy()

// 	// create
// 	ensure.Nil(t, db.Put(wo, givenKey, givenVal1))

// 	// retrieve
// 	v1, err := db.Get(ro, givenKey)
// 	defer v1.Free()
// 	ensure.Nil(t, err)
// 	ensure.DeepEqual(t, v1.Data(), givenVal1)

// 	// retrieve bytes
// 	_v1, err := db.GetBytes(ro, givenKey)
// 	ensure.Nil(t, err)
// 	ensure.DeepEqual(t, _v1, givenVal1)

// 	// update
// 	ensure.Nil(t, db.Put(wo, givenKey, givenVal2))
// 	v2, err := db.Get(ro, givenKey)
// 	defer v2.Free()
// 	ensure.Nil(t, err)
// 	ensure.DeepEqual(t, v2.Data(), givenVal2)

// 	// retrieve pinned
// 	v3, err := db.GetPinned(ro, givenKey)
// 	defer v3.Destroy()
// 	ensure.Nil(t, err)
// 	ensure.DeepEqual(t, v3.Data(), givenVal2)

// 	engine, err := CreateBackupEngine(db)
// 	require.Nil(t, err)
// 	defer engine.Close()

// 	t.Run("createBackupAndVerify", func(t *testing.T) {
// 		infos := engine.GetInfo()
// 		require.Empty(t, infos)

// 		// create first backup
// 		require.Nil(t, engine.CreateNewBackup())

// 		// create second backup
// 		require.Nil(t, engine.CreateNewBackupFlush(true))

// 		infos = engine.GetInfo()
// 		require.Equal(t, 2, len(infos))
// 		for i := range infos {
// 			require.Nil(t, engine.VerifyBackup(infos[i].ID))
// 			require.True(t, infos[i].Size > 0)
// 			require.True(t, infos[i].NumFiles > 0)
// 		}
// 	})

// 	t.Run("purge", func(t *testing.T) {
// 		require.Nil(t, engine.PurgeOldBackups(1))

// 		infos := engine.GetInfo()
// 		require.Equal(t, 1, len(infos))
// 	})

// 	t.Run("restoreFromLatest", func(t *testing.T) {
// 		dir, err := ioutil.TempDir("", "gorocksdb-restoreFromLatest")
// 		require.Nil(t, err)

// 		ro := NewRestoreOptions()
// 		defer ro.Destroy()
// 		require.Nil(t, engine.RestoreDBFromLatestBackup(dir, dir, ro))
// 		require.Nil(t, engine.RestoreDBFromLatestBackup(dir, dir, ro))
// 	})

// 	t.Run("restoreFromBackup", func(t *testing.T) {
// 		infos := engine.GetInfo()
// 		require.Equal(t, 1, len(infos))

// 		dir, err := ioutil.TempDir("", "gorocksdb-restoreFromBackup")
// 		require.Nil(t, err)

// 		ro := NewRestoreOptions()
// 		defer ro.Destroy()
// 		require.Nil(t, engine.RestoreDBFromBackup(dir, dir, ro, infos[0].ID))

// 		// try to reopen restored db
// 		backupDB, err := OpenDb(db.opts, dir)
// 		require.Nil(t, err)

// 		r := NewDefaultReadOptions()
// 		defer r.Destroy()

// 		v3, err := backupDB.GetPinned(r, givenKey)
// 		defer v3.Destroy()
// 		ensure.Nil(t, err)
// 		ensure.DeepEqual(t, v3.Data(), givenVal2)
// 	})
// }
