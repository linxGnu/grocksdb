FROM centos:7

# install toolchain(s) and go
ENV GOLANG_PACKAGE go1.13.4.linux-amd64.tar.gz

RUN yum -y update && \
    yum -y install gcc gcc-c++ git pkg-config make which unzip && \
    curl https://dl.google.com/go/${GOLANG_PACKAGE} -o ${GOLANG_PACKAGE} && \
    tar -C /usr/local -xzf ${GOLANG_PACKAGE} && rm ${GOLANG_PACKAGE} && \
    yum clean all && rm -rf /var/cache/yum

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN mkdir -p "$GOPATH/src/github.com/linxGnu/grocksdb" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# install cmake
RUN yum install -y wget && \
    cd /tmp && \
    wget https://github.com/Kitware/CMake/releases/download/v3.15.5/cmake-3.15.5.tar.gz && \
    tar xzf cmake-3.15.5.tar.gz && cd cmake-3.15.5 && \
    ./bootstrap --parallel=16 && make -j16 && make install && \
    cd /tmp && rm -rf * && \
    yum remove -y wget && yum clean all && rm -rf /var/cache/yum

# building
ADD . "$GOPATH/src/github.com/linxGnu/grocksdb"
