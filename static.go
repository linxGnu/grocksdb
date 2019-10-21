// +build static

package grocksdb

// #cgo LDFLAGS: -lrocksdb -lstdc++ -lm -ldl -lzstd -llz4 -lz -lsnappy -lbz2
import "C"
