package grocksdb

// #include "rocksdb/c.h"
import "C"

// Cache is a cache used to store data read from data in memory.
type Cache struct {
	c *C.rocksdb_cache_t
}

// NewLRUCache creates a new LRU Cache object with the capacity given.
func NewLRUCache(capacity uint64) *Cache {
	return NewNativeCache(C.rocksdb_cache_create_lru(C.size_t(capacity)))
}

// NewNativeCache creates a Cache object.
func NewNativeCache(c *C.rocksdb_cache_t) *Cache {
	return &Cache{c}
}

// GetUsage returns the Cache memory usage.
func (c *Cache) GetUsage() uint64 {
	return uint64(C.rocksdb_cache_get_usage(c.c))
}

// GetPinnedUsage returns the Cache pinned memory usage.
func (c *Cache) GetPinnedUsage() uint64 {
	return uint64(C.rocksdb_cache_get_pinned_usage(c.c))
}

// SetCapacity sets capacity of the cache.
func (c *Cache) SetCapacity(value uint64) {
	C.rocksdb_cache_set_capacity(c.c, C.size_t(value))
}

// // GetCapacity returns capacity of the cache.
// func (c *Cache) GetCapacity() uint64 {
// 	return uint64(C.rocksdb_cache_get_capacity(c.c))
// }

// Destroy deallocates the Cache object.
func (c *Cache) Destroy() {
	C.rocksdb_cache_destroy(c.c)
	c.c = nil
}
