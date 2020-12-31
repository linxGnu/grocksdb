package grocksdb

import "testing"

func TestOptEnv(t *testing.T) {
	opt := NewDefaultEnvOptions()
	defer opt.Destroy()
}
