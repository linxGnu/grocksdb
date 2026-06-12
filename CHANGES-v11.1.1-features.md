# grocksdb — New feature wrappers for RocksDB v10.10.1 → v11.1.1

This document records the **Go-level wrappers added** for RocksDB C API surface that
landed between RocksDB **v10.10.1** and **v11.1.1**.

## Context

grocksdb is a thin **CGO** binding over RocksDB's C API (`c.h`). The vendored
repo-root `c.h` had already been synced to v11.1.1, so most of the new C symbols were
**declared but not wrapped** in Go. A reviewer pointed out that the upgrade PR was
missing useful wrappers (their example: `rocksdb_block_based_options_set_index_block_search_type`).

This change adds idiomatic Go wrappers for the new, usefully-wrappable C API,
following grocksdb's existing patterns (enum/bool setters, native-handle factories
with `Destroy()`, and `//export` + `COWList` cgo callbacks).

## What was added

### 1. Block-based table options — `options_block_based_table.go`

| Go API | C symbol |
|--------|----------|
| `IndexBlockSearchType` enum (`KBinarySearchIndexBlockSearchType`, `KInterpolationSearchIndexBlockSearchType`) | enum `rocksdb_block_based_table_index_block_search_type_*` |
| `(*BlockBasedTableOptions) SetIndexBlockSearchType(IndexBlockSearchType)` | `rocksdb_block_based_options_set_index_block_search_type` |
| `(*BlockBasedTableOptions) SetSeparateKeyValueInDataBlock(bool)` | `rocksdb_block_based_options_set_separate_key_value_in_data_block` |
| `(*BlockBasedTableOptions) SetBlockAlign(bool)` | `rocksdb_block_based_options_set_block_align` |

### 2. FIFO compaction options — `options_compaction.go`

| Go API | C symbol |
|--------|----------|
| `(*FIFOCompactionOptions) Set/GetMaxDataFilesSize(uint64)` | `rocksdb_fifo_compaction_options_{set,get}_max_data_files_size` |
| `(*FIFOCompactionOptions) SetUseKVRatioCompaction(bool)` / `UseKVRatioCompaction() bool` | `rocksdb_fifo_compaction_options_{set,get}_use_kv_ratio_compaction` |

### 3. DB-level options & factories — `options.go`, `file_checksum.go`, `sst_partitioner.go`

| Go API | C symbol |
|--------|----------|
| `(*Options) SetOpenFilesAsync(bool)` / `OpenFilesAsync() bool` | `rocksdb_options_{set,get}_open_files_async` |
| `NewCRC32CFileChecksumGenFactory()` + `(*FileChecksumGenFactory) Destroy()` | `rocksdb_file_checksum_gen_crc32c_factory_create` / `..._factory_destroy` |
| `(*Options) SetFileChecksumGenFactory(*FileChecksumGenFactory)` | `rocksdb_options_set_file_checksum_gen_factory` |
| `NewFixedPrefixSSTPartitionerFactory(uint)` + `(*SSTPartitionerFactory) Destroy()` | `rocksdb_sst_partitioner_fixed_prefix_factory_create` / `..._factory_destroy` |
| `(*Options) SetSSTPartitionerFactory(*SSTPartitionerFactory)` | `rocksdb_options_set_sst_partitioner_factory` |

> **Ownership:** `SetFileChecksumGenFactory` / `SetSSTPartitionerFactory` take
> ownership of the factory's native handle. The `Options` frees it on `Destroy()`;
> callers should **not** also call `Destroy()` on the passed factory. New
> `Options` struct fields `cfcg` / `cspf` hold the handles and are freed in
> `(*Options).Destroy()`.

### 4. Callback logger — `logger.go`, `grocksdb.h`, `grocksdb.c`

| Go API | C symbol |
|--------|----------|
| `LogHandler` type `func(level InfoLogLevel, msg string)` | — |
| `NewCallbackLogger(InfoLogLevel, LogHandler) *Logger` | `rocksdb_logger_create_callback_logger` (via `gorocksdb_logger_create_callback` shim) |

Implemented with the comparator-style cgo pattern: a `COWList` registry holds the
handler, the registry index is passed as the callback's `priv`, and
`//export gorocksdb_logger_log` forwards each message back to Go. Reuses the existing
`(*Logger).Destroy()` and `(*Options).SetInfoLog`.

### 5. WriteBatch lazy-data iteration — `write_batch.go`, `grocksdb.h`, `grocksdb.c`

| Go API | C symbol |
|--------|----------|
| `WriteBatchHandler` interface (`Put`, `Delete`, `LogData`, `PutCF`, `DeleteCF`, `MergeCF`) | — |
| `(*WriteBatch) Iterate(WriteBatchHandler)` | `rocksdb_writebatch_iterate_ld` (via `gorocksdb_writebatch_iterate_ld`) |
| `(*WriteBatch) IterateCF(WriteBatchHandler)` | `rocksdb_writebatch_iterate_cf_ld` (via `gorocksdb_writebatch_iterate_cf_ld`) |

The `_ld` (lazy-data) variants add a `log_data` callback over the existing
iterate/iterate_cf signatures and handle large values. A `COWList` registry plus six
`//export` callbacks (`gorocksdb_writebatch_put`, `_deleted`, `_logdata`, `_put_cf`,
`_deleted_cf`, `_merge_cf`) bridge to the handler. The existing pure-Go
`(*WriteBatch) NewIterator()` is left unchanged for backward compatibility.

## Files changed

**Modified:** `options_block_based_table.go`, `options_compaction.go`, `options.go`,
`logger.go`, `write_batch.go`, `grocksdb.h`, `grocksdb.c`.

**Added:** `file_checksum.go`, `sst_partitioner.go`.

**Tests:** extended `options_block_based_table_test.go`, `options_compaction_test.go`,
`options_test.go` (`TestOptionsOpenFilesAsyncAndFactories`),
`write_batch_test.go` (`TestWriteBatchIterateLD`); added `logger_test.go`
(`TestCallbackLogger`).

## Intentionally out of scope

| Feature | Reason |
|---------|--------|
| `rocksdb_abort_all_compactions` / `rocksdb_resume_all_compactions` | Not present in the vendored `c.h` yet; would require a `c.h` edit first. |
| `rocksdb_table_properties_collector_factory_*` | `c.h` exposes only `_destroy` and `_add` — there is **no C constructor**, so nothing can be created to add. |
| `rocksdb_compaction_service_options_override_*` (~17 fns) | Large, advanced remote-compaction surface; deferred. |

## Validation

- `gofmt` clean on all changed/added files.
- cgo cross-check: the 7 `//export` Go callbacks match the 7 C trampoline call sites.
- Pending (requires the cgo toolchain + linked RocksDB v11.1.1):
  ```
  go build ./...
  go vet ./...
  go test -run 'TestBBT|TestFifoCompactOption|TestOptions|CallbackLogger|WriteBatchIterateLD' ./...
  ```
