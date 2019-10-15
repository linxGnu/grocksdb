GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)

DEST=dist/$(GOOS_GOARCH)
DEST_INCLUDE=$(DEST)/include

MAKEFLAGS = -j4
CFLAGS += ${EXTRA_CFLAGS}
CXXFLAGS += ${EXTRA_CXXFLAGS}
LDFLAGS += $(EXTRA_LDFLAGS)
MACHINE ?= $(shell uname -m)
ARFLAGS = ${EXTRA_ARFLAGS} rs
STRIPFLAGS = -S -x

BZIP2_VER ?= 1.0.6
BZIP2_SHA256 ?= a2848f34fcd5d6cf47def00461fcb528a0484d8edef8208d6d2e2909dc61d9cd
BZIP2_DOWNLOAD_BASE ?= https://web.archive.org/web/20180624184835/http://www.bzip.org
SHA256_CMD = sha256sum

default: prepare zlib bz2 snappy lz4 zstd rocksdb

.PHONY: prepare
prepare:
	mkdir -p $(DEST_INCLUDE)/zlib $(DEST_INCLUDE)/bz2 $(DEST_INCLUDE)/snappy $(DEST_INCLUDE)/lz4 $(DEST_INCLUDE)/zstd

.PHONY: zlib
zlib:
	git submodule update --remote --init --recursive -- libs/zlib
	cd libs/zlib && CFLAGS='-fPIC ${EXTRA_CFLAGS}' LDFLAGS='${EXTRA_LDFLAGS}' ./configure --static && $(MAKE) $(MAKEFLAGS)
	cp libs/zlib/libz.a $(DEST)/
	cp libs/zlib/*.h $(DEST_INCLUDE)/zlib/

.PHONY: bz2
bz2:
	pushd libs ; \
	curl --output bzip2-$(BZIP2_VER).tar.gz -L ${BZIP2_DOWNLOAD_BASE}/$(BZIP2_VER)/bzip2-$(BZIP2_VER).tar.gz ; \
	BZIP2_SHA256_ACTUAL=`$(SHA256_CMD) bzip2-$(BZIP2_VER).tar.gz | cut -d ' ' -f 1`; \
	if [ "$(BZIP2_SHA256)" != "$$BZIP2_SHA256_ACTUAL" ]; then \
		echo bzip2-$(BZIP2_VER).tar.gz checksum mismatch, expected=\"$(BZIP2_SHA256)\" actual=\"$$BZIP2_SHA256_ACTUAL\"; \
		exit 1; \
	fi ; \
	tar xzf bzip2-$(BZIP2_VER).tar.gz && rm bzip2-$(BZIP2_VER).tar.gz && mv bzip2-$(BZIP2_VER) bzip2 ; \
	cd bzip2 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 -g -D_FILE_OFFSET_BITS=64 ${EXTRA_CFLAGS}' AR='ar ${EXTRA_ARFLAGS}' libbz2.a ; \
	popd ; \
	cp libs/bzip2/libbz2.a $(DEST)/
	cp libs/bzip2/*.h $(DEST_INCLUDE)/bz2/

.PHONY: snappy
snappy:
	git submodule update --remote --init --recursive -- libs/snappy
	mkdir -p libs/snappy/build && cd libs/snappy/build && CFLAGS='${EXTRA_CFLAGS}' CXXFLAGS='${EXTRA_CXXFLAGS}' LDFLAGS='${EXTRA_LDFLAGS}' cmake -DCMAKE_POSITION_INDEPENDENT_CODE=ON .. && $(MAKE) $(MAKEFLAGS)
	cp libs/snappy/build/libsnappy.a $(DEST)/
	cp libs/snappy/*.h $(DEST_INCLUDE)/snappy/

.PHONY: lz4
lz4:
	git submodule update --remote --init --recursive -- libs/lz4
	cd libs/lz4 && $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' all
	cp libs/lz4/lib/liblz4.a $(DEST)/
	cp libs/lz4/lib/*.h $(DEST_INCLUDE)/lz4/

.PHONY: zstd
zstd:
	git submodule update --remote --init --recursive -- libs/zstd
	cd libs/zstd/lib && DESTDIR=. PREFIX= $(MAKE) $(MAKEFLAGS) CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' install
	cp libs/zstd/lib/libzstd.a $(DEST)/
	cp libs/zstd/lib/include/*.h $(DEST_INCLUDE)/zstd/

.PHONY: rocksdb
rocksdb:
	git submodule update --remote --init --recursive -- libs/rocksdb
	cd libs/rocksdb && CFLAGS='-fPIC -O2 ${EXTRA_CFLAGS}' CXXFLAGS='-fPIC -O2 ${EXTRA_CXXFLAGS}' $(MAKE) $(MAKEFLAGS) static_lib
	cd libs/rocksdb && strip $(STRIPFLAGS) librocksdb.a
	cp libs/rocksdb/librocksdb.a $(DEST)/
	cp -R libs/rocksdb/include/rocksdb $(DEST_INCLUDE)/