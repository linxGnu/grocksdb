// +build static !builtin_static

package grocksdb

// #cgo LDFLAGS: -lrocksdb -pthread -ldl -lm -lzstd -llz4 -lz -lsnappy
import "C"
