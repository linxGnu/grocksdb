#!/bin/bash

# update sys
yum update -y

# install tool
yum install -y gcc gcc-c++ git pkg-config make which

# install deps (for linking) and temporary build tool
yum install -y snappy-devel cmake

# build rocksdb
make

# remove deps
yum remove -y snappy-devel cmake
