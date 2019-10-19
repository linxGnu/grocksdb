#!/bin/bash

# update sys
yum update -y

# install deps
yum install -y wget which make gcc gcc-c++ snappy-devel

# install cmake
wget https://cmake.org/files/v3.11/cmake-3.11.4.tar.gz
tar xzf cmake-3.11.4.tar.gz
pushd cmake-3.11.4
./bootstrap --prefix=/opt/cmake --no-system-libs --parallel=16
make install -j16
popd
rm -rf cmake-3.11.4 cmake-3.11.4.tar.gz

# build rocksdb
PATH=$PATH:/opt/cmake/bin
make

# remove snappy-devel
yum remove -y snappy-devel

# test
go test -v -tags static_rocksdb
