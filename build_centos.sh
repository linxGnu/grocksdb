#!/bin/bash

# update sys
yum update -y

# install tool
yum install -y gcc gcc-c++ git pkg-config make which

# install temporary build tool
yum install -y cmake

# build rocksdb
make

# remove temporary build tool
yum remove -y cmake
