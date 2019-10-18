// +build static_rocksdb

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include
// #cgo LDFLAGS: -L ${SRCDIR}/dist/linux_amd64/lib -lrocksdb -lstdc++ -lm -ldl -lzstd -llz4 -lz -lsnappy -lbz2
import "C"
