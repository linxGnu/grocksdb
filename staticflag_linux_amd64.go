// +build static

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/libs/linux_amd64/include
// #cgo CXXFLAGS: -I ${SRCDIR}/libs/linux_amd64/include -fno-rtti -std=gnu++11
// #cgo LDFLAGS: -L ${SRCDIR}/libs/linux_amd64 -lz -llz4 -lzstd -lsnappy -lrocksdb -lm -ldl -lstdc++
import "C"
