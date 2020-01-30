// +build static

package grocksdb

// #cgo LDFLAGS: -static -lrocksdb -pthread -lstdc++ -ldl -lm -lzstd -llz4 -lz -lsnappy -lgflags
import "C"
