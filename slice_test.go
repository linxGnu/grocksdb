package grocksdb

import "testing"

func TestStringToSlice(t *testing.T) {
	slice := StringToSlice("asdf")
	defer slice.Free()
}
