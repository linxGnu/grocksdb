#!/bin/sh

# install deps for linking
brew install bzip2 lz4 snappy zlib

# build rocksdb
make deps
make rocksdb

# remove deps
brew remove bzip2 lz4 snappy zlib

# test
go test -v -tags static_rocksdb
