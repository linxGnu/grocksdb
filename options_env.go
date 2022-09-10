package grocksdb

// #include "rocksdb/c.h"
import "C"
import "unsafe"

// EnvOptions represents options for env.
type EnvOptions struct {
	c *C.rocksdb_envoptions_t
}

// NewDefaultEnvOptions creates a default EnvOptions object.
func NewDefaultEnvOptions() *EnvOptions {
	return NewNativeEnvOptions(C.rocksdb_envoptions_create())
}

// NewNativeEnvOptions creates a EnvOptions object.
func NewNativeEnvOptions(c *C.rocksdb_envoptions_t) *EnvOptions {
	return &EnvOptions{c: c}
}

// Destroy deallocates the EnvOptions object.
func (opts *EnvOptions) Destroy() {
	C.rocksdb_envoptions_destroy(opts.c)
	opts.c = nil
}

// Native returns native EnvOptions
func (opts *EnvOptions) Native() unsafe.Pointer {
	return unsafe.Pointer(opts.c)
}
