#!/bin/bash

# update sys
yum update -y

# install gcc-8
yum install -y centos-release-scl
yum install -y devtoolset-8-gcc devtoolset-8-gcc-c++
scl enable devtoolset-8 -- bash

# install tool
yum install -y which git

# install deps
yum install -y snappy-devel

# install cmake
curl https://cmake.org/files/v3.11/cmake-3.11.4.tar.gz -o cmake-3.11.4.tar.gz
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

# remove cmake
rm -rf /opt/cmake
