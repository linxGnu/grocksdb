// +build !linux !static

package grocksdb

// #cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -lsnappy -ldl
import "C"
