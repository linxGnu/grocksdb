#!/bin/sh

# install snappy for linking
brew install snappy

# build rocksdb
make

# remove snappy
brew remove snappy
