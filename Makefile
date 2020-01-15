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
SNAPPY_COMMIT = 537f4ad6240e586970fe554614542e9717df7902
LZ4_COMMIT = 0f749838bf29bc0d1df428e23cf3dbb76ec4e9fc
ZSTD_COMMIT = 10f0e6993f9d2f682da6d04aa2385b7d53cbb4ee
BZ2_COMMIT = 6a8690fc8d26c815e798c588f796eabe9d684cf0
ROCKSDB_COMMIT = f48aa1c3084700bc72b73ee36f027e428f0dda86
JEMALLOC_COMMIT = ea6b3e973b477b8061e0076bb257dbd7f3faa756

ROCKSDB_EXTRA_CXXFLAGS := 
ifeq ($(GOOS), darwin)
	ROCKSDB_EXTRA_CXXFLAGS += -fPIC -O3 -w -I$(DEST_INCLUDE)
else
	ROCKSDB_EXTRA_CXXFLAGS += -fPIC -O3 -Wno-error=shadow -I$(DEST_INCLUDE)
endif

default: prepare jemalloc zlib snappy bz2 lz4 zstd rocksdb

.PHONY: prepare
prepare:
	rm -rf $(DEST)
	mkdir -p $(DEST_LIB) $(DEST_INCLUDE)

.PHONY: jemalloc
jemalloc:
	git submodule update --remote --init --recursive -- libs/jemalloc
	cd libs/jemalloc && git checkout $(JEMALLOC_COMMIT)
	cd libs/jemalloc && sh autogen.sh && CFLAGS='-fPIC -O3 ${EXTRA_CFLAGS}' ./configure --prefix=$(DEST) --enable-prof && make build_lib_static && make install_lib_static install_include

.PHONY: zlib
zlib:
	git submodule update --remote --init --recursive -- libs/zlib
	cd libs/zlib && git checkout $(ZLIB_COMMIT)
	cd libs/zlib && CFLAGS='-fPIC -O3 ${EXTRA_CFLAGS}' ./configure --static && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) all
	cp libs/zlib/libz.a $(DEST_LIB)/
	cp libs/zlib/*.h $(DEST_INCLUDE)/

.PHONY: snappy
snappy:
	git submodule update --remote --init --recursive -- libs/snappy
	cd libs/snappy && git checkout $(SNAPPY_COMMIT)
	cd libs/snappy && rm -rf build && mkdir -p build && cd build && \
	CFLAGS='-O3 ${EXTRA_CFLAGS}' cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON .. && \
	$(MAKE) clean && $(MAKE) $(MAKE_FLAGS) snappy
	cp libs/snappy/build/libsnappy.a $(DEST_LIB)/
	cp libs/snappy/build/*.h $(DEST_INCLUDE)/
	cp libs/snappy/*.h $(DEST_INCLUDE)/

.PHONY: lz4
lz4:
	git submodule update --remote --init --recursive -- libs/lz4
	cd libs/lz4 && git checkout $(LZ4_COMMIT)
	cd libs/lz4 && $(MAKE) clean && $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O3 ${EXTRA_CFLAGS}' lz4-release
	cp libs/lz4/lib/liblz4.a $(DEST_LIB)/
	cp libs/lz4/lib/*.h $(DEST_INCLUDE)/

.PHONY: zstd
zstd:
	git submodule update --remote --init --recursive -- libs/zstd
	cd libs/zstd && git checkout $(ZSTD_COMMIT)
	cd libs/zstd/lib && $(MAKE) clean && DESTDIR=. PREFIX= $(MAKE) $(MAKE_FLAGS) CFLAGS='-fPIC -O3 ${EXTRA_CFLAGS}' default install
	cp libs/zstd/lib/libzstd.a $(DEST_LIB)/
	cp libs/zstd/lib/include/*.h $(DEST_INCLUDE)/

.PHONY: bz2
bz2:
	cd libs/bzip2 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O3 -g -D_FILE_OFFSET_BITS=64 ${EXTRA_CFLAGS}' AR='ar ${EXTRA_ARFLAGS}' bzip2
	cp libs/bzip2/libbz2.a $(DEST_LIB)/
	cp libs/bzip2/*.h $(DEST_INCLUDE)/

.PHONY: rocksdb
rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb
	cd libs/rocksdb && git checkout $(ROCKSDB_COMMIT) && mkdir -p build && cd build && cmake -DCMAKE_LIBRARY_PATH=${DEST}/lib -DCMAKE_INCLUDE_PATH=${DEST}/include \
	-DWITH_BZ2=1  -DCMAKE_BUILD_TYPE=Release -DWITH_JEMALLOC=1 -DWITH_SNAPPY=1 -DWITH_LZ4=1 -DWITH_ZLIB=1 -DWITH_ZSTD=1 .. && \
	EXTRA_CXXFLAGS='$(ROCKSDB_EXTRA_CXXFLAGS)' EXTRA_LDFLAGS='-L$(DEST_LIB)' $(MAKE) $(MAKE_FLAGS) rocksdb
	cd libs/rocksdb/build && strip $(STRIPFLAGS) librocksdb.a
	cp libs/rocksdb/build/librocksdb.a $(DEST_LIB)/
	cp -R libs/rocksdb/include/rocksdb $(DEST_INCLUDE)/

.PHONY: test
test:
	go test -v -count=1 -tags builtin_static
