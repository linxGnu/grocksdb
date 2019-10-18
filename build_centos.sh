#!/bin/bash

# update sys
yum update -y

# install deps
yum install -y snappy-devel zlib-devel bzip2-devel lz4-devel wget which make gcc gcc-c++

# install cmake
wget https://cmake.org/files/v3.11/cmake-3.11.4.tar.gz
tar xzf cmake-3.11.4.tar.gz
pushd cmake-3.11.4
./bootstrap --prefix=/opt/cmake --no-system-libs
make install -j8
popd
rm -rf cmake-3.11.4 cmake-3.11.4.tar.gz

# install zstd
wget https://github.com/facebook/zstd/releases/download/v1.4.3/zstd-1.4.3.tar.gz
tar xzf zstd-1.4.3.tar.gz
pushd zstd-1.4.3
make install -j8
popd
rm -rf zstd-1.4.3 zstd-1.4.3.tar.gz

# build rocksdb
PATH=$PATH:/opt/cmake/bin
make all

# remove deps
yum remove -y snappy-devel zlib-devel bzip2-devel lz4-devel

# test
go test -v -tags static_rocksdb
