package grocksdb

// #include "rocksdb/c.h"
import "C"

// SSTFileManager is used to track SST and blob files in the DB and control
// their deletion rate. All SSTFileManager public functions are thread-safe.
type SSTFileManager struct {
	c *C.rocksdb_sst_file_manager_t
}

// NewSSTFileManager creates new SSTFileManager.
func NewSSTFileManager(env *Env) *SSTFileManager {
	return &SSTFileManager{c: C.rocksdb_sst_file_manager_create(env.c)}
}

// Destroy deallocates the SSTFileManager object.
func (s *SSTFileManager) Destroy() {
	C.rocksdb_sst_file_manager_destroy(s.c)
	s.c = nil
}

// SetMaxAllowedSpaceUsage updates the maximum allowed space that should be used by RocksDB, if
// the total size of the SST and blob files exceeds max_allowed_space, writes
// to RocksDB will fail.
//
// Setting max_allowed_space to 0 will disable this feature; maximum allowed
// space will be infinite (Default value).
//
// thread-safe.
func (s *SSTFileManager) SetMaxAllowedSpaceUsage(space uint64) {
	C.rocksdb_sst_file_manager_set_max_allowed_space_usage(s.c, C.uint64_t(space))
}

// SetCompactionBufferSize sets the amount of buffer room each compaction should be able to leave.
// In other words, at its maximum disk space consumption, the compaction
// should still leave compaction_buffer_size available on the disk so that
// other background functions may continue, such as logging and flushing.
//
// thread-safe.
func (s *SSTFileManager) SetCompactionBufferSize(size uint64) {
	C.rocksdb_sst_file_manager_set_compaction_buffer_size(s.c, C.uint64_t(size))
}

// IsMaxAllowedSpaceReached returns true if the total size of SST  and blob files exceeded the maximum
// allowed space usage.
//
// thread-safe.
func (s *SSTFileManager) IsMaxAllowedSpaceReached() bool {
	return bool(C.rocksdb_sst_file_manager_is_max_allowed_space_reached(s.c))
}

// IsMaxAllowedSpaceReachedIncludingCompactions returns true if the total size of SST and blob files as well as estimated
// size of ongoing compactions exceeds the maximums allowed space usage.
func (s *SSTFileManager) IsMaxAllowedSpaceReachedIncludingCompactions() bool {
	return bool(C.rocksdb_sst_file_manager_is_max_allowed_space_reached_including_compactions(s.c))
}

// GetTotalSize returns the total size of all tracked files.
//
// thread-safe
func (s *SSTFileManager) GetTotalSize() uint64 {
	return uint64(C.rocksdb_sst_file_manager_get_total_size(s.c))
}

// GetDeleteRateBytesPerSecond returns delete rate limit in bytes per second.
//
// thread-safe
func (s *SSTFileManager) GetDeleteRateBytesPerSecond() int64 {
	return int64(C.rocksdb_sst_file_manager_get_delete_rate_bytes_per_second(s.c))
}

// SetDeleteRateBytesPerSecond updates the delete rate limit in bytes per second.
// zero means disable delete rate limiting and delete files immediately.
//
// thread-safe
func (s *SSTFileManager) SetDeleteRateBytesPerSecond(rate int64) {
	C.rocksdb_sst_file_manager_set_delete_rate_bytes_per_second(s.c, C.int64_t(rate))
}

// GetMaxTrashDBRatio returns trash/DB size ratio where new files will be deleted immediately.
//
// thread-safe
func (s *SSTFileManager) GetMaxTrashDBRatio() float64 {
	return float64(C.rocksdb_sst_file_manager_get_max_trash_db_ratio(s.c))
}

// SetMaxTrashDBRatio updates trash/DB size ratio where new files will be deleted immediately.
//
// thread-safe
func (s *SSTFileManager) SetMaxTrashDBRatio(ratio float64) {
	C.rocksdb_sst_file_manager_set_max_trash_db_ratio(s.c, C.double(ratio))
}

// GetTotalTrashSize returns the total size of trash files.
//
// thread-safe
func (s *SSTFileManager) GetTotalTrashSize() uint64 {
	return uint64(C.rocksdb_sst_file_manager_get_total_trash_size(s.c))
}
