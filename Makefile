GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOOS_GOARCH := $(GOOS)_$(GOARCH)
GOOS_GOARCH_NATIVE := $(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)
LIBROCKSDB_NAME := librocksdb_$(GOOS_GOARCH).a

.PHONY: librocksdb.a

librocksdb.a: $(LIBROCKSDB_NAME)

clean:
	rm -f $(LIBROCKSDB_NAME)
	rm -rf libs/rocksdb

static:
	git submodule update --remote --init --recursive -- libs/rocksdb ; \
	pushd libs/rocksdb; \
	make -j4 static_lib ; \
	popd ; \
	sh repack.sh libs/rocksdb/librocksdb.a ; \
	mv librocksdb.a $(LIBROCKSDB_NAME)
