#include "grocksdb.h"
#include "_cgo_export.h"

/* Base */

void gorocksdb_destruct_handler(void* state) { }

/* Comparator */

rocksdb_comparator_t* gorocksdb_comparator_create(uintptr_t idx) {
    return rocksdb_comparator_create(
        (void*)idx,
        gorocksdb_destruct_handler,
        (int (*)(void*, const char*, size_t, const char*, size_t))(gorocksdb_comparator_compare),
        (const char *(*)(void*))(gorocksdb_comparator_name));
}

rocksdb_comparator_t* gorocksdb_comparator_with_ts_create(uintptr_t idx, size_t ts_size) {
    return rocksdb_comparator_with_ts_create(
        (void*)idx,
        gorocksdb_destruct_handler,
        (int (*)(void*, const char*, size_t, const char*, size_t))(gorocksdb_comparator_compare),
        (int (*)(void*, const char*, size_t, const char*, size_t))(gorocksdb_comparator_compare_ts),
        (int (*)(void*, const char*, size_t, unsigned char, const char*, size_t, unsigned char))(gorocksdb_comparator_compare_without_ts),
        (const char* (*)(void*))(gorocksdb_comparator_name),
        ts_size);
}

/* CompactionFilter */

rocksdb_compactionfilter_t* gorocksdb_compactionfilter_create(uintptr_t idx) {
    return rocksdb_compactionfilter_create(
        (void*)idx,
        gorocksdb_destruct_handler,
        (unsigned char (*)(void*, int, const char*, size_t, const char*, size_t, char**, size_t*, unsigned char*))(gorocksdb_compactionfilter_filter),
        (const char *(*)(void*))(gorocksdb_compactionfilter_name));
}

/* Merge Operator */

rocksdb_mergeoperator_t* gorocksdb_mergeoperator_create(uintptr_t idx) {
    return rocksdb_mergeoperator_create(
        (void*)idx,
        gorocksdb_destruct_handler,
        (char* (*)(void*, const char*, size_t, const char*, size_t, const char* const*, const size_t*, int, unsigned char*, size_t*))(gorocksdb_mergeoperator_full_merge),
        (char* (*)(void*, const char*, size_t, const char* const*, const size_t*, int, unsigned char*, size_t*))(gorocksdb_mergeoperator_partial_merge_multi),
        gorocksdb_mergeoperator_delete_value,
        (const char* (*)(void*))(gorocksdb_mergeoperator_name));
}

void gorocksdb_mergeoperator_delete_value(void* id, const char* v, size_t s) {
    free((char*)v);
}

/* Slice Transform */

rocksdb_slicetransform_t* gorocksdb_slicetransform_create(uintptr_t idx) {
    return rocksdb_slicetransform_create(
    	(void*)idx,
    	gorocksdb_destruct_handler,
    	(char* (*)(void*, const char*, size_t, size_t*))(gorocksdb_slicetransform_transform),
    	(unsigned char (*)(void*, const char*, size_t))(gorocksdb_slicetransform_in_domain),
    	(const char* (*)(void*))(gorocksdb_slicetransform_name));
}

/* Logger */

/* Trampoline: RocksDB invokes this with priv == the registry index, and we
   forward to the exported Go callback. */
static void gorocksdb_logger_logv(void* priv, unsigned lev, char* msg, size_t len) {
    gorocksdb_logger_log((uintptr_t)priv, lev, msg, len);
}

rocksdb_logger_t* gorocksdb_logger_create_callback(int log_level, uintptr_t idx) {
    return rocksdb_logger_create_callback_logger(log_level, gorocksdb_logger_logv, (void*)idx);
}

/* WriteBatch lazy-data iteration */

/* Non-CF trampolines forward state (the registry index) to the exported Go
   callbacks. */
static void gorocksdb_wb_put(void* s, const char* k, size_t klen, const char* v, size_t vlen) {
    gorocksdb_writebatch_put((uintptr_t)s, (char*)k, klen, (char*)v, vlen);
}
static void gorocksdb_wb_deleted(void* s, const char* k, size_t klen) {
    gorocksdb_writebatch_deleted((uintptr_t)s, (char*)k, klen);
}
static void gorocksdb_wb_logdata(void* s, const char* blob, size_t blen) {
    gorocksdb_writebatch_logdata((uintptr_t)s, (char*)blob, blen);
}

void gorocksdb_writebatch_iterate_ld(rocksdb_writebatch_t* wb, uintptr_t idx) {
    rocksdb_writebatch_iterate_ld(wb, (void*)idx,
        gorocksdb_wb_put, gorocksdb_wb_deleted, gorocksdb_wb_logdata);
}

/* CF trampolines. */
static void gorocksdb_wb_put_cf(void* s, uint32_t cfid, const char* k, size_t klen, const char* v, size_t vlen) {
    gorocksdb_writebatch_put_cf((uintptr_t)s, cfid, (char*)k, klen, (char*)v, vlen);
}
static void gorocksdb_wb_deleted_cf(void* s, uint32_t cfid, const char* k, size_t klen) {
    gorocksdb_writebatch_deleted_cf((uintptr_t)s, cfid, (char*)k, klen);
}
static void gorocksdb_wb_merge_cf(void* s, uint32_t cfid, const char* k, size_t klen, const char* v, size_t vlen) {
    gorocksdb_writebatch_merge_cf((uintptr_t)s, cfid, (char*)k, klen, (char*)v, vlen);
}

void gorocksdb_writebatch_iterate_cf_ld(rocksdb_writebatch_t* wb, uintptr_t idx) {
    rocksdb_writebatch_iterate_cf_ld(wb, (void*)idx,
        gorocksdb_wb_put_cf, gorocksdb_wb_deleted_cf, gorocksdb_wb_merge_cf,
        gorocksdb_wb_logdata);
}
