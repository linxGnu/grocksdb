// +build !linux !static

package grocksdb

// #cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -ldl
import "C"
