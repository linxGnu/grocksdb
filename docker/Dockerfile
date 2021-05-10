FROM centos:7

# install toolchain
RUN yum -y update && yum -y --setopt=tsflags=nodocs install gcc gcc-c++ git pkg-config make which unzip automake wget libtool && \
    yum clean all && rm -rf /var/cache/yum

# openssl dev 
ARG ssl_version="OpenSSL_1_1_1k"
RUN cd /tmp && wget https://github.com/openssl/openssl/archive/${ssl_version}.tar.gz && tar xzf ${ssl_version}.tar.gz && cd openssl-${ssl_version}  && \
    ./Configure linux-x86_64 -Wa,--noexecstack no-shared no-dso -DDSO_NONE --prefix=/usr/local --openssldir=/usr/local && make -j$(nproc) install_dev && \
    cd /tmp && rm -rf * && ldconfig

# install latest cmake
ARG cmake_version="3.20.2"
RUN cd /tmp && wget https://github.com/Kitware/CMake/releases/download/v${cmake_version}/cmake-${cmake_version}.tar.gz && tar xzf cmake-${cmake_version}.tar.gz && \
    cd cmake-${cmake_version} && ./bootstrap --parallel=$(nproc) && make -j$(nproc) && make install && \
    cd /tmp && rm -rf *

# setup cxx standard
ENV CFLAGS='-fPIC -O3 -pipe'
ENV CXXFLAGS='-fPIC -O3 -pipe'

# install go
ENV GOLANG_PACKAGE go1.16.4.linux-amd64.tar.gz

RUN curl https://dl.google.com/go/${GOLANG_PACKAGE} -o ${GOLANG_PACKAGE} && \
    tar -C /usr/local -xzf ${GOLANG_PACKAGE} && rm ${GOLANG_PACKAGE}

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# cleanup
RUN yum remove -y automake libtool && yum clean all && rm -rf /var/cache/yum

# setup pkg-config path
ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:/usr/local/lib64/pkgconfig
