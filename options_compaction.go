package grocksdb

// #include "rocksdb/c.h"
import "C"

// UniversalCompactionStopStyle describes a algorithm used to make a
// compaction request stop picking new files into a single compaction run.
type UniversalCompactionStopStyle uint

// Compaction stop style types.
const (
	CompactionStopStyleSimilarSize = UniversalCompactionStopStyle(C.rocksdb_similar_size_compaction_stop_style)
	CompactionStopStyleTotalSize   = UniversalCompactionStopStyle(C.rocksdb_total_size_compaction_stop_style)
)

// BottommostLevelCompaction for level based compaction, we can configure if we want to skip/force
// bottommost level compaction.
type BottommostLevelCompaction byte

const (
	// KSkip skip bottommost level compaction
	KSkip BottommostLevelCompaction = 0
	// KIfHaveCompactionFilter only compact bottommost level if there is a compaction filter
	// This is the default option
	KIfHaveCompactionFilter BottommostLevelCompaction = 1
	// KForce always compact bottommost level
	KForce BottommostLevelCompaction = 2
	// KForceOptimized always compact bottommost level but in bottommost level avoid
	// double-compacting files created in the same compaction
	KForceOptimized BottommostLevelCompaction = 3
)

// CompactRangeOptions represent all of the available options for compact range.
type CompactRangeOptions struct {
	c *C.rocksdb_compactoptions_t
}

// NewCompactRangeOptions creates new compact range options.
func NewCompactRangeOptions() *CompactRangeOptions {
	return &CompactRangeOptions{
		c: C.rocksdb_compactoptions_create(),
	}
}

// Destroy deallocates the CompactionOptions object.
func (opts *CompactRangeOptions) Destroy() {
	C.rocksdb_compactoptions_destroy(opts.c)
	opts.c = nil
}

// SetExclusiveManualCompaction if more than one thread calls manual compaction,
// only one will actually schedule it while the other threads will simply wait
// for the scheduled manual compaction to complete. If exclusive_manual_compaction
// is set to true, the call will disable scheduling of automatic compaction jobs
// and wait for existing automatic compaction jobs to finish.
func (opts *CompactRangeOptions) SetExclusiveManualCompaction(value bool) {
	C.rocksdb_compactoptions_set_exclusive_manual_compaction(opts.c, boolToChar(value))
}

// SetBottommostLevelCompaction set bottommost level compaction.
func (opts *CompactRangeOptions) SetBottommostLevelCompaction(value BottommostLevelCompaction) {
	C.rocksdb_compactoptions_set_bottommost_level_compaction(opts.c, C.uchar(value))
}

// SetChangeLevel if true, compacted files will be moved to the minimum level capable
// of holding the data or given level (specified non-negative target_level).
func (opts *CompactRangeOptions) SetChangeLevel(value bool) {
	C.rocksdb_compactoptions_set_change_level(opts.c, boolToChar(value))
}

// SetTargetLevel if change_level is true and target_level have non-negative value, compacted
// files will be moved to target_level.
func (opts *CompactRangeOptions) SetTargetLevel(value int32) {
	C.rocksdb_compactoptions_set_target_level(opts.c, C.int(value))
}

// FIFOCompactionOptions represent all of the available options for
// FIFO compaction.
type FIFOCompactionOptions struct {
	c *C.rocksdb_fifo_compaction_options_t
}

// NewDefaultFIFOCompactionOptions creates a default FIFOCompactionOptions object.
func NewDefaultFIFOCompactionOptions() *FIFOCompactionOptions {
	return NewNativeFIFOCompactionOptions(C.rocksdb_fifo_compaction_options_create())
}

// NewNativeFIFOCompactionOptions creates a native FIFOCompactionOptions object.
func NewNativeFIFOCompactionOptions(c *C.rocksdb_fifo_compaction_options_t) *FIFOCompactionOptions {
	return &FIFOCompactionOptions{c}
}

// SetMaxTableFilesSize sets the max table file size.
// Once the total sum of table files reaches this, we will delete the oldest
// table file
// Default: 1GB
func (opts *FIFOCompactionOptions) SetMaxTableFilesSize(value uint64) {
	C.rocksdb_fifo_compaction_options_set_max_table_files_size(opts.c, C.uint64_t(value))
}

// Destroy deallocates the FIFOCompactionOptions object.
func (opts *FIFOCompactionOptions) Destroy() {
	C.rocksdb_fifo_compaction_options_destroy(opts.c)
}

// UniversalCompactionOptions represent all of the available options for
// universal compaction.
type UniversalCompactionOptions struct {
	c *C.rocksdb_universal_compaction_options_t
}

// NewDefaultUniversalCompactionOptions creates a default UniversalCompactionOptions
// object.
func NewDefaultUniversalCompactionOptions() *UniversalCompactionOptions {
	return NewNativeUniversalCompactionOptions(C.rocksdb_universal_compaction_options_create())
}

// NewNativeUniversalCompactionOptions creates a UniversalCompactionOptions
// object.
func NewNativeUniversalCompactionOptions(c *C.rocksdb_universal_compaction_options_t) *UniversalCompactionOptions {
	return &UniversalCompactionOptions{c}
}

// SetSizeRatio sets the percentage flexibilty while comparing file size.
// If the candidate file(s) size is 1% smaller than the next file's size,
// then include next file into this candidate set.
// Default: 1
func (opts *UniversalCompactionOptions) SetSizeRatio(value uint) {
	C.rocksdb_universal_compaction_options_set_size_ratio(opts.c, C.int(value))
}

// SetMinMergeWidth sets the minimum number of files in a single compaction run.
// Default: 2
func (opts *UniversalCompactionOptions) SetMinMergeWidth(value uint) {
	C.rocksdb_universal_compaction_options_set_min_merge_width(opts.c, C.int(value))
}

// SetMaxMergeWidth sets the maximum number of files in a single compaction run.
// Default: UINT_MAX
func (opts *UniversalCompactionOptions) SetMaxMergeWidth(value uint) {
	C.rocksdb_universal_compaction_options_set_max_merge_width(opts.c, C.int(value))
}

// SetMaxSizeAmplificationPercent sets the size amplification.
// It is defined as the amount (in percentage) of
// additional storage needed to store a single byte of data in the database.
// For example, a size amplification of 2% means that a database that
// contains 100 bytes of user-data may occupy upto 102 bytes of
// physical storage. By this definition, a fully compacted database has
// a size amplification of 0%. Rocksdb uses the following heuristic
// to calculate size amplification: it assumes that all files excluding
// the earliest file contribute to the size amplification.
// Default: 200, which means that a 100 byte database could require upto
// 300 bytes of storage.
func (opts *UniversalCompactionOptions) SetMaxSizeAmplificationPercent(value uint) {
	C.rocksdb_universal_compaction_options_set_max_size_amplification_percent(opts.c, C.int(value))
}

// SetCompressionSizePercent sets the percentage of compression size.
//
// If this option is set to be -1, all the output files
// will follow compression type specified.
//
// If this option is not negative, we will try to make sure compressed
// size is just above this value. In normal cases, at least this percentage
// of data will be compressed.
// When we are compacting to a new file, here is the criteria whether
// it needs to be compressed: assuming here are the list of files sorted
// by generation time:
//    A1...An B1...Bm C1...Ct
// where A1 is the newest and Ct is the oldest, and we are going to compact
// B1...Bm, we calculate the total size of all the files as total_size, as
// well as  the total size of C1...Ct as total_C, the compaction output file
// will be compressed iff
//   total_C / total_size < this percentage
// Default: -1
func (opts *UniversalCompactionOptions) SetCompressionSizePercent(value int) {
	C.rocksdb_universal_compaction_options_set_compression_size_percent(opts.c, C.int(value))
}

// SetStopStyle sets the algorithm used to stop picking files into a single compaction run.
// Default: CompactionStopStyleTotalSize
func (opts *UniversalCompactionOptions) SetStopStyle(value UniversalCompactionStopStyle) {
	C.rocksdb_universal_compaction_options_set_stop_style(opts.c, C.int(value))
}

// Destroy deallocates the UniversalCompactionOptions object.
func (opts *UniversalCompactionOptions) Destroy() {
	C.rocksdb_universal_compaction_options_destroy(opts.c)
	opts.c = nil
}
