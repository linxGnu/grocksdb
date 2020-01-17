// +build builtin_static

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/dist/darwin_amd64/include
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/darwin_amd64/include
// #cgo LDFLAGS: -L ${SRCDIR}/dist/darwin_amd64/lib -pthread -lrocksdb -lstdc++ -lm -ldl -lzstd -llz4 -lz -lsnappy -lbz2
import "C"
