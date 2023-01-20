//go:build !testing && !grocksdb_clean_link

package grocksdb

// #cgo LDFLAGS: -lrocksdb -pthread -lstdc++ -ldl -lm -lzstd -llz4 -lz -lsnappy
import "C"
