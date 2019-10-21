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

### Default - Builtin Static (Linux, Mac OS)

`grocksdb` contains built static version of `Rocksdb` with:
- gcc (GCC) 4.8.5 20150623 (Red Hat 4.8.5-39)
- Apple clang version 11.0.0 (clang-1100.0.33.8). 

You have to do nothing on your machine. Just install it like other go libraries:

```bash
go get -u github.com/linxGnu/grocksdb

# Build your project with `builtin_static` tags:
go build -tags builtin_static
```

### Static lib (Linux, Mac OS)

If you don't trust my builtin/want to build with your compiler/env:

##### Prerequisite
- cmake 3.11+
- make

##### Build

Make sure to install libraries for linking before making targets.

```bash
# You could find `Makefile` at root of repository

# build static libs
make
```

Then, build your project with tags (same as above):

```
go build -tags builtin_static
```

### Existed Static lib

In case, already have static-lib version of rocksdb, you could build your project with:

```
# don't use builtin static
go build -tags static
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
