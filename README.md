# grocksdb, a Go wrapper for RocksDB

[![Build Status](https://travis-ci.org/linxGnu/grocksdb.svg?branch=master)](https://travis-ci.org/linxGnu/grocksdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/grocksdb)](https://goreportcard.com/report/github.com/linxGnu/grocksdb)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/grocksdb)

This is a `Fork` from [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb). I respect the author work and community contribution.
The `LICENSE` still remains as upstream.

Why I made a patched clone instead of PR:
- Supports almost C API (unlike upstream). Catching up with latest version of Rocksdb as promise.
- Static build focused.
- This fork contains `no defer` in codebase (my side project requires as less overhead as possible). This introduces loose
convention of how/when to free c-mem, thus break the rule of [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb).

## Install

### Default (Linux, Mac OS)

`grocksdb` contains built static version of `Rocksdb`. You have to do nothing on your machine. Just install it like other go libraries:

```bash
go get -u github.com/linxGnu/grocksdb

# Build your project with `static_rocksdb` tags:
go build -tags static_rocksdb
```

### Static lib (Linux, Mac OS)

If you don't trust my built-ready static version, you could build your own:

#### For CentOS/Mac OS:

```bash
# Below scripts could be found in root of repo.

# Centos
sh build_centos.sh

# Mac OS
sh build_macos.sh
```

#### Others:

##### Prerequisite
- cmake 3.11+
- make
- For linking purpose:
  - libsnappy-dev/snappy-devel
  - zlib1g-dev/zlib-devel
  - libbz2-dev/bzip2-devel
  - liblz4-dev/lz4-devel

##### Build

Make sure to install libraries for linking before making targets.

```bash
# build static libs
make deps
make rocksdb

# On success, you could remove above linked libraries safely
```

### Shared lib

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v6.3.6+ on your machine.

After that, you can install `gorocksdb` using the following command:

    CGO_CFLAGS="-I/path/to/rocksdb/include" \
    CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
      go get github.com/linxGnu/grocksdb

## Usage

See also: [doc](https://godoc.org/github.com/linxGnu/grocksdb)

## API Support

Almost C API, excepts:
- [ ] putv/mergev/deletev/delete_rangev
- [ ] compaction_filter_factory/compaction_filter_context
