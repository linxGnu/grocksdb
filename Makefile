GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)

ROOT_DIR=${PWD}
DEST=$(ROOT_DIR)/dist/$(GOOS_GOARCH)
DEST_INCLUDE=$(DEST)/include

MAKE_FLAGS = -j8
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
ZSTD_COMMIT = 8b6d96827c24dd09109830272f413254833317d9
BZ2_COMMIT = 6a8690fc8d26c815e798c588f796eabe9d684cf0
ROCKSDB_COMMIT = d47cdbc1888440a75ecf43646fd1ddab8ebae9be

deps: prepare zlib snappy bz2 lz4 zstd

default_target: deps rocksdb

.PHONY: prepare
prepare:
	rm -rf $(DEST)
	mkdir -p $(DEST_INCLUDE)/zlib $(DEST_INCLUDE)/snappy $(DEST_INCLUDE)/lz4 $(DEST_INCLUDE)/zstd $(DEST_INCLUDE)/bz2

.PHONY: zlib
zlib:
	git submodule update --remote --init --recursive -- libs/zlib
	cd libs/zlib && git checkout $(ZLIB_COMMIT)
	cd libs/zlib && CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' LDFLAGS='${EXTRA_LDFLAGS}' ./configure --static && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) all
	cp libs/zlib/libz.a $(DEST)/
	cp libs/zlib/*.h $(DEST_INCLUDE)/zlib/

.PHONY: snappy
snappy:
	git submodule update --remote --init --recursive -- libs/snappy
	cd libs/snappy && git checkout $(SNAPPY_COMMIT)
	cd libs/snappy && rm -rf build && mkdir -p build && cd build && \
	CFLAGS='-O2 ${EXTRA_CFLAGS}' CXXFLAGS='-O2 ${EXTRA_CXXFLAGS}' LDFLAGS='${EXTRA_LDFLAGS}' cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON .. && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) snappy
	cp libs/snappy/build/libsnappy.a $(DEST)/
	cp libs/snappy/*.h $(DEST_INCLUDE)/snappy/

.PHONY: lz4
lz4:
	git submodule update --remote --init --recursive -- libs/lz4
	cd libs/lz4 && git checkout $(LZ4_COMMIT)
	cd libs/lz4 && $(MAKE) clean && $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' lz4 lz4-release
	cp libs/lz4/lib/liblz4.a $(DEST)/
	cp libs/lz4/lib/*.h $(DEST_INCLUDE)/lz4/

.PHONY: zstd
zstd:
	git submodule update --remote --init --recursive -- libs/zstd
	cd libs/zstd && git checkout $(ZSTD_COMMIT)
	cd libs/zstd/lib && $(MAKE) clean && DESTDIR=. PREFIX= $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' all install
	cp libs/zstd/lib/libzstd.a $(DEST)/
	cp libs/zstd/lib/include/*.h $(DEST_INCLUDE)/zstd/

.PHONY: bz2
bz2:
	git submodule update --remote --init --recursive -- libs/bz2
	cd libs/bz2 && git checkout $(BZ2_COMMIT)
	cd libs/bz2 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 -g -D_FILE_OFFSET_BITS=64 ${EXTRA_CFLAGS}' AR='ar ${EXTRA_ARFLAGS}' bzip2
	cp libs/bz2/libbz2.a $(DEST)/
	cp libs/bz2/*.h $(DEST_INCLUDE)/bz2/

.PHONY: rocksdb
rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb
	cd libs/rocksdb && git checkout $(ROCKSDB_COMMIT) && $(MAKE) clean && \
	CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' CXXFLAGS='-fPIC -O2 ${EXTRA_CXXFLAGS} -Wno-error=shadow' \
	USE_RTTI=1 $(MAKE) $(MAKE_FLAGS) static_lib
	cd libs/rocksdb && strip $(STRIPFLAGS) librocksdb.a
	cp libs/rocksdb/librocksdb.a $(DEST)/
	cp -R libs/rocksdb/include/rocksdb $(DEST_INCLUDE)/
