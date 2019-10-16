#!/bin/bash

# update sys
yum update -y

# install deps
yum install -y snappy-devel zlib-devel bzip2-devel lz4-devel wget git which make

# install cmake
wget https://cmake.org/files/v3.11/cmake-3.11.4.tar.gz
tar zxvf cmake-3.11.4.tar.gz
./bootstrap --prefix=/opt/cmake --no-system-libs
make install -j8
