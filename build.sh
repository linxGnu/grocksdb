#!/bin/bash
DIRECTORY="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

INSTALL_PREFIX=$1

export CFLAGS='-fPIC -O3 -pipe -funroll-loops' 
export CXXFLAGS='-fPIC -O3 -pipe -funroll-loops'

BUILD_PATH=/tmp/build
mkdir -p $BUILD_PATH

CMAKE_REQUIRED_PARAMS="-DCMAKE_POSITION_INDEPENDENT_CODE=ON -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}"

cd $BUILD_PATH && wget https://github.com/madler/zlib/archive/v1.2.11.tar.gz && tar xzf v1.2.11.tar.gz && cd zlib-1.2.11 && \
    ./configure --prefix=$INSTALL_PREFIX --static && make -j16 install && \
    cd $BUILD_PATH && rm -rf *

cd $BUILD_PATH && wget https://github.com/google/snappy/archive/1.1.8.tar.gz && tar xzf 1.1.8.tar.gz && cd snappy-1.1.8 && \
    mkdir -p build_place && cd build_place && cmake $CMAKE_REQUIRED_PARAMS -DSNAPPY_BUILD_TESTS=OFF .. && make install/strip -j16 && \
    cd $BUILD_PATH && rm -rf *

cd $BUILD_PATH && wget https://github.com/lz4/lz4/archive/v1.9.2.tar.gz && tar xzf v1.9.2.tar.gz && cd lz4-1.9.2/contrib/cmake_unofficial && \
    cmake $CMAKE_REQUIRED_PARAMS -DLZ4_BUILD_LEGACY_LZ4C=OFF -DBUILD_SHARED_LIBS=OFF -DLZ4_POSITION_INDEPENDENT_LIB=ON && make -j16 install && \
    cd $BUILD_PATH && rm -rf *

cd $BUILD_PATH && wget https://github.com/facebook/zstd/releases/download/v1.4.4/zstd-1.4.4.tar.gz && tar xzf zstd-1.4.4.tar.gz && cd zstd-1.4.4/build/cmake && mkdir -p build_place && cd build_place && \
    cmake -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX} -DZSTD_BUILD_PROGRAMS=OFF -DZSTD_BUILD_CONTRIB=OFF -DZSTD_BUILD_STATIC=ON -DZSTD_BUILD_SHARED=OFF -DZSTD_BUILD_TESTS=OFF \
    $CMAKE_REQUIRED_PARAMS -DZSTD_ZLIB_SUPPORT=ON -DZSTD_LZMA_SUPPORT=OFF -DCMAKE_BUILD_TYPE=Release .. && make -j16 install && \
    cd $BUILD_PATH && rm -rf *

cd $BUILD_PATH && wget https://github.com/facebook/rocksdb/archive/v6.6.3.tar.gz && tar xzf v6.6.3.tar.gz && cd rocksdb-6.6.3/ && \
    cp $DIRECTORY/CMakeLists.txt ./ && \
    mkdir -p build_place && cd build_place && cmake -DCMAKE_BUILD_TYPE=Release $CMAKE_REQUIRED_PARAMS -DCMAKE_PREFIX_PATH=$INSTALL_PREFIX -DWITH_TESTS=OFF -DWITH_BENCHMARK_TOOLS=OFF -DWITH_TOOLS=OFF \
    -DWITH_MD_LIBRARY=OFF -DWITH_RUNTIME_DEBUG=OFF -DROCKSDB_BUILD_SHARED=OFF -DWITH_SNAPPY=ON -DWITH_LZ4=ON -DWITH_ZLIB=ON -DWITH_ZSTD=ON -DWITH_BZ2=OFF -WITH_GFLAGS=OFF .. && make -j16 install/strip && \
    cd $BUILD_PATH && rm -rf *

rm -rf $INSTALL_PREFIX/bin $INSTALL_PREFIX/share $INSTALL_PREFIX/lib/cmake $INSTALL_PREFIX/lib64/cmake $INSTALL_PREFIX/lib/pkgconfig $INSTALL_PREFIX/lib64/pkgconfig
