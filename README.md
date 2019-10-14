# grocksdb, a Go wrapper for RocksDB

[![Build Status](https://travis-ci.org/linxGnu/grocksdb.svg?branch=master)](https://travis-ci.org/linxGnu/grocksdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/grocksdb)](https://goreportcard.com/report/github.com/linxGnu/grocksdb)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/grocksdb)

This is a `Fork` from [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb). I respect the author work and community contribution.
The `LICENSE` still remains as upstream.

Why I made a patched clone instead of PR:
- Unlike upstream, this fork contains `no defer` in codebase (my side project requires as less overhead as possible). This introduces loose
convention of how/when to free c-mem, thus break the rule of [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb).
- Catching up with latest version of Rocksdb.

## Install

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v6.3.6+ on your machine.

After that, you can install gorocksdb using the following command:

    CGO_CFLAGS="-I/path/to/rocksdb/include" \
    CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
      go get github.com/linxGnu/grocksdb

## Support

Almost C API, excepts:
- [ ] putv/mergev/deletev/delete_rangev
- [ ] compaction_filter_factory/compaction_filter_context
- [ ] writebatch_iterate
