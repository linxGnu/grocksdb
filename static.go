// +build static

package grocksdb

// #cgo CFLAGS: -static
// #cgo CXXFLAGS: -static
// #cgo LDFLAGS: -lrocksdb -pthread -lstdc++ -ldl -lm -lzstd -llz4 -lz -lsnappy -lgflags
import "C"
