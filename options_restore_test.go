package grocksdb

import "testing"

func TestRestoreOption(t *testing.T) {
	ro := NewRestoreOptions()
	defer ro.Destroy()

	ro.SetKeepLogFiles(123)
}
