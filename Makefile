GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)
ROOT_DIR=${PWD}

clean:
	rm -rf libs/* ; \
	mkdir -p libs/$(GOOS_GOARCH)/include ; \

lz4:
	git submodule update --remote --init --recursive -- libs/lz4 ; \
	pushd libs/lz4 ; \
	CXXFLAGS="-fno-rtti" make -j4 lib ; \
	popd ; \
	cp libs/lz4/lib/*.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/lz4/lib/liblz4.a libs/$(GOOS_GOARCH)/

zlib:
	git submodule update --remote --init --recursive -- libs/zlib ; \
	pushd libs/zlib ; \
	./configure ; \
	CXXFLAGS="-fno-rtti" make -j4 all ; \
	popd ; \
	cp libs/zlib/*.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/zlib/libz.a libs/$(GOOS_GOARCH)/

snappy:
	git submodule update --remote --init --recursive -- libs/snappy ; \
	pushd libs/snappy ; \
	cmake . && CXXFLAGS="-fno-rtti" make -j4 ; \
	popd ; \
	cp libs/snappy/*.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/snappy/libsnappy.a libs/$(GOOS_GOARCH)/

zstd:
	git submodule update --remote --init --recursive -- libs/zstd ; \
	pushd libs/zstd ; \
	make -j4 lib-release ; \
	popd ; \
	cp libs/zstd/lib/zstd.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/zstd/lib/common/zstd_errors.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/zstd/lib/deprecated/zbuff.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/zstd/lib/dictBuilder/zdict.h libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/zstd/lib/libzstd.a libs/$(GOOS_GOARCH)/

deps: clean lz4 zlib snappy zstd

rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb ; \
	pushd libs/rocksdb ; \
	LDFLAGS="-L$(ROOT_DIR)/libs/$(GOOS_GOARCH) -lm -lstdc++ -lz -llz4 -lsnappy -lzstd" \
	CFLAGS="-I$(ROOT_DIR)/libs/$(GOOS_GOARCH)/include" CXXFLAGS="-I$(ROOT_DIR)/libs/$(GOOS_GOARCH)/include" \
	USE_RTTI=1 make -j4 static_lib ; \
	popd ; \
	cp -R libs/rocksdb/include/rocksdb libs/$(GOOS_GOARCH)/include/ ; \
	cp libs/rocksdb/librocksdb.a libs/$(GOOS_GOARCH)/ ; \
	pushd libs/$(GOOS_GOARCH) ; \
	sh ../../repack.sh librocksdb.a ; \
	popd
