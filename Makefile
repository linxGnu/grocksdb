GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)

ROOT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DEST=$(ROOT_DIR)/dist/$(GOOS_GOARCH)
DEST_LIB=$(DEST)/lib
DEST_INCLUDE=$(DEST)/include

MAKE_FLAGS = -j16
CFLAGS += ${EXTRA_CFLAGS}
CXXFLAGS += ${EXTRA_CXXFLAGS}
LDFLAGS += $(EXTRA_LDFLAGS)
MACHINE ?= $(shell uname -m)
ARFLAGS = ${EXTRA_ARFLAGS} rs
STRIPFLAGS = -S -x

# Dependencies and Rocksdb
ZLIB_COMMIT = cacf7f1d4e3d44d871b605da3b647f07d718623f
SNAPPY_COMMIT = e9e11b84e629c3e06fbaa4f0a86de02ceb9d6992
LZ4_COMMIT = e8baeca51ef2003d6c9ec21c32f1563fef1065b9
ZSTD_COMMIT = ed65210c9b6635e21e67e60138f86d04a071681f
BZ2_COMMIT = 6a8690fc8d26c815e798c588f796eabe9d684cf0
ROCKSDB_COMMIT = d47cdbc1888440a75ecf43646fd1ddab8ebae9be

default: prepare zlib snappy bz2 lz4 zstd rocksdb

.PHONY: prepare
prepare:
	rm -rf $(DEST)
	mkdir -p $(DEST_LIB) $(DEST_INCLUDE)

.PHONY: zlib
zlib:
	git submodule update --remote --init --recursive -- libs/zlib
	cd libs/zlib && git checkout $(ZLIB_COMMIT)
	cd libs/zlib && CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' ./configure --static && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) all
	cp libs/zlib/libz.a $(DEST_LIB)/
	cp libs/zlib/*.h $(DEST_INCLUDE)/
	cp libs/zlib/libz.a /usr/local/lib/
	cp libs/zlib/*.h /usr/local/include/

.PHONY: snappy
snappy:
	git submodule update --remote --init --recursive -- libs/snappy
	cd libs/snappy && git checkout $(SNAPPY_COMMIT)
	cd libs/snappy && rm -rf build && mkdir -p build && cd build && \
	CFLAGS='-O2 ${EXTRA_CFLAGS}' cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON .. && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) install
	cp libs/snappy/build/libsnappy.a $(DEST_LIB)/
	cp libs/snappy/*.h $(DEST_INCLUDE)/

.PHONY: lz4
lz4:
	git submodule update --remote --init --recursive -- libs/lz4
	cd libs/lz4 && git checkout $(LZ4_COMMIT)
	cd libs/lz4 && $(MAKE) clean && $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' lz4 lz4-release
	cp libs/lz4/lib/liblz4.a $(DEST_LIB)/
	cp libs/lz4/lib/*.h $(DEST_INCLUDE)/
	cp libs/lz4/lib/liblz4.a /usr/local/lib/
	cp libs/lz4/lib/*.h /usr/local/include/

.PHONY: zstd
zstd:
	git submodule update --remote --init --recursive -- libs/zstd
	cd libs/zstd && git checkout $(ZSTD_COMMIT)
	cd libs/zstd/lib && $(MAKE) clean && DESTDIR=. PREFIX= $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' default install
	cp libs/zstd/lib/libzstd.a $(DEST_LIB)/
	cp libs/zstd/lib/include/*.h $(DEST_INCLUDE)/
	cp libs/zstd/lib/libzstd.a /usr/local/lib/
	cp libs/zstd/lib/include/*.h /usr/local/include/

.PHONY: bz2
bz2:
	cd libs/bzip2 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 -g -D_FILE_OFFSET_BITS=64 ${EXTRA_CFLAGS}' AR='ar ${EXTRA_ARFLAGS}' bzip2
	cp libs/bzip2/libbz2.a $(DEST_LIB)/
	cp libs/bzip2/*.h $(DEST_INCLUDE)/
	cp libs/bzip2/libbz2.a /usr/local/lib/
	cp libs/bzip2/*.h /usr/local/include/

.PHONY: rocksdb
rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb
	cd libs/rocksdb && git checkout $(ROCKSDB_COMMIT) && $(MAKE) clean && \
	CXXFLAGS='-fPIC -O2 -Wno-error=shadow ${EXTRA_CXXFLAGS}' $(MAKE) $(MAKE_FLAGS) static_lib
	cd libs/rocksdb && strip $(STRIPFLAGS) librocksdb.a
	cp libs/rocksdb/librocksdb.a $(DEST_LIB)/
	cp -R libs/rocksdb/include/rocksdb $(DEST_INCLUDE)/
