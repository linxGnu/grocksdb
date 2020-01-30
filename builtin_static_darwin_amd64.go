// +build builtin_static !static

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/dist/darwin_amd64/include
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/darwin_amd64/include
// #cgo LDFLAGS: -L${SRCDIR}/dist/darwin_amd64/lib -lrocksdb -pthread -lstdc++ -ldl -lm -lzstd -llz4 -lz -lsnappy
import "C"
