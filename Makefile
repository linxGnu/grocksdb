GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)

ROOT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DEST=$(ROOT_DIR)/dist/$(GOOS_GOARCH)
DEST_LIB=$(DEST)/lib
DEST_INCLUDE=$(DEST)/include

MAKE_FLAGS = -j16
MACHINE ?= $(shell uname -m)
STRIPFLAGS = -S -x

# Dependencies and Rocksdb
ZLIB_COMMIT = cacf7f1d4e3d44d871b605da3b647f07d718623f
SNAPPY_COMMIT = e9e11b84e629c3e06fbaa4f0a86de02ceb9d6992
LZ4_COMMIT = e8baeca51ef2003d6c9ec21c32f1563fef1065b9
ZSTD_COMMIT = a3d655d2255481333e09ecca9855f1b37f757c52
BZ2_COMMIT = 6a8690fc8d26c815e798c588f796eabe9d684cf0
ROCKSDB_COMMIT = e3169e3ea8762d2f34880742106858a23c8dc8b7

ROCKSDB_EXTRA_CXXFLAGS := 
ifeq ($(GOOS), darwin)
	ROCKSDB_EXTRA_CXXFLAGS += -fPIC -O2 -w -I$(DEST_INCLUDE) -DZLIB -DBZIP2 -DSNAPPY -DLZ4 -DZSTD
else
	ROCKSDB_EXTRA_CXXFLAGS += -fPIC -O2 -Wno-error=shadow -I$(DEST_INCLUDE) -DZLIB -DBZIP2 -DSNAPPY -DLZ4 -DZSTD
endif

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

.PHONY: snappy
snappy:
	git submodule update --remote --init --recursive -- libs/snappy
	cd libs/snappy && git checkout $(SNAPPY_COMMIT)
	cd libs/snappy && rm -rf build && mkdir -p build && cd build && \
	CFLAGS='-O2 ${EXTRA_CFLAGS}' cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON .. && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) snappy
	cp libs/snappy/build/libsnappy.a $(DEST_LIB)/
	cp libs/snappy/build/*.h $(DEST_INCLUDE)/
	cp libs/snappy/*.h $(DEST_INCLUDE)/

.PHONY: lz4
lz4:
	git submodule update --remote --init --recursive -- libs/lz4
	cd libs/lz4 && git checkout $(LZ4_COMMIT)
	cd libs/lz4 && $(MAKE) clean && $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' lz4-release
	cp libs/lz4/lib/liblz4.a $(DEST_LIB)/
	cp libs/lz4/lib/*.h $(DEST_INCLUDE)/

.PHONY: zstd
zstd:
	git submodule update --remote --init --recursive -- libs/zstd
	cd libs/zstd && git checkout $(ZSTD_COMMIT)
	cd libs/zstd/lib && $(MAKE) clean && DESTDIR=. PREFIX= $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' default install
	cp libs/zstd/lib/libzstd.a $(DEST_LIB)/
	cp libs/zstd/lib/include/*.h $(DEST_INCLUDE)/

.PHONY: bz2
bz2:
	cd libs/bzip2 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 -g -D_FILE_OFFSET_BITS=64 ${EXTRA_CFLAGS}' AR='ar ${EXTRA_ARFLAGS}' bzip2
	cp libs/bzip2/libbz2.a $(DEST_LIB)/
	cp libs/bzip2/*.h $(DEST_INCLUDE)/

.PHONY: rocksdb
rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb
	cd libs/rocksdb && git checkout $(ROCKSDB_COMMIT) && $(MAKE) clean && \
	$(MAKE) $(MAKE_FLAGS) EXTRA_CXXFLAGS='$(ROCKSDB_EXTRA_CXXFLAGS)' EXTRA_LDFLAGS='-L$(DEST_LIB)' static_lib
	cd libs/rocksdb && strip $(STRIPFLAGS) librocksdb.a
	cp libs/rocksdb/librocksdb.a $(DEST_LIB)/
	cp -R libs/rocksdb/include/rocksdb $(DEST_INCLUDE)/

.PHONY: test
test:
	go test -v -count=1 -tags builtin_static
