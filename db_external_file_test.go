package grocksdb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalFile(t *testing.T) {
	db := newTestDB(t, "TestDBExternalFile", nil)
	defer db.Close()

	envOpts := NewDefaultEnvOptions()
	opts := NewDefaultOptions()
	w := NewSSTFileWriter(envOpts, opts)
	defer w.Destroy()

	filePath, err := ioutil.TempFile("", "sst-file-test")
	require.Nil(t, err)
	defer os.Remove(filePath.Name())

	err = w.Open(filePath.Name())
	require.Nil(t, err)

	err = w.Add([]byte("aaa"), []byte("aaaValue"))
	require.Nil(t, err)
	err = w.Add([]byte("bbb"), []byte("bbbValue"))
	require.Nil(t, err)
	err = w.Add([]byte("ccc"), []byte("cccValue"))
	require.Nil(t, err)
	err = w.Add([]byte("ddd"), []byte("dddValue"))
	require.Nil(t, err)

	err = w.Finish()
	require.Nil(t, err)

	ingestOpts := NewDefaultIngestExternalFileOptions()
	err = db.IngestExternalFile([]string{filePath.Name()}, ingestOpts)
	require.Nil(t, err)

	readOpts := NewDefaultReadOptions()

	v1, err := db.Get(readOpts, []byte("aaa"))
	require.Nil(t, err)
	require.EqualValues(t, v1.Data(), []byte("aaaValue"))
	v2, err := db.Get(readOpts, []byte("bbb"))
	require.Nil(t, err)
	require.EqualValues(t, v2.Data(), []byte("bbbValue"))
	v3, err := db.Get(readOpts, []byte("ccc"))
	require.Nil(t, err)
	require.EqualValues(t, v3.Data(), []byte("cccValue"))
	v4, err := db.Get(readOpts, []byte("ddd"))
	require.Nil(t, err)
	require.EqualValues(t, v4.Data(), []byte("dddValue"))
}
