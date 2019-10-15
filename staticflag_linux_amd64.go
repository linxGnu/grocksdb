// +build static

package grocksdb

// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include
// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/bz2
// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/lz4
// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/snappy
// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/zlib
// #cgo CFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/zstd
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/bz2
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/lz4
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/snappy
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/zlib
// #cgo CXXFLAGS: -I ${SRCDIR}/dist/linux_amd64/include/zstd
// #cgo LDFLAGS: -L ${SRCDIR}/dist/linux_amd64 -lz -llz4 -lzstd -lsnappy -lrocksdb -lm -ldl -lstdc++
import "C"
