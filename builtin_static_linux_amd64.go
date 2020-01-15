// +build builtin_static

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include -fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free
// #cgo LDFLAGS: -L ${SRCDIR}/dist/linux_amd64/lib -lrocksdb -lstdc++ -lm -ldl -lzstd -llz4 -lz -lsnappy -lbz2 -ljemalloc_pic
import "C"
