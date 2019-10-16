package grocksdb

// #include "rocksdb/c.h"
import "C"

// Env is a system call environment used by a database.
type Env struct {
	c *C.rocksdb_env_t
}

// NewDefaultEnv creates a default environment.
func NewDefaultEnv() *Env {
	return NewNativeEnv(C.rocksdb_create_default_env())
}

// NewMemEnv returns a new environment that stores its data in memory and delegates
// all non-file-storage tasks to base_env. The caller must delete the result
// when it is no longer needed.
// *base_env must remain live while the result is in use.
func NewMemEnv() *Env {
	return NewNativeEnv(C.rocksdb_create_mem_env())
}

// NewNativeEnv creates a Environment object.
func NewNativeEnv(c *C.rocksdb_env_t) *Env {
	return &Env{c}
}

// SetBackgroundThreads sets the number of background worker threads
// of a specific thread pool for this environment.
// 'LOW' is the default pool.
// Default: 1
func (env *Env) SetBackgroundThreads(n int) {
	C.rocksdb_env_set_background_threads(env.c, C.int(n))
}

// SetHighPriorityBackgroundThreads sets the size of the high priority
// thread pool that can be used to prevent compactions from stalling
// memtable flushes.
func (env *Env) SetHighPriorityBackgroundThreads(n int) {
	C.rocksdb_env_set_high_priority_background_threads(env.c, C.int(n))
}

// JoinAllThreads wait for all threads started by StartThread to terminate.
func (env *Env) JoinAllThreads() {
	C.rocksdb_env_join_all_threads(env.c)
}

// LowerThreadPoolIOPriority lower IO priority for threads from the specified pool.
func (env *Env) LowerThreadPoolIOPriority() {
	C.rocksdb_env_lower_thread_pool_io_priority(env.c)
}

// LowerHighPriorityThreadPoolIOPriority lower IO priority for high priority
// thread pool.
func (env *Env) LowerHighPriorityThreadPoolIOPriority() {
	C.rocksdb_env_lower_high_priority_thread_pool_io_priority(env.c)
}

// LowerThreadPoolCPUPriority lower CPU priority for threads from the specified pool.
func (env *Env) LowerThreadPoolCPUPriority() {
	C.rocksdb_env_lower_thread_pool_cpu_priority(env.c)
}

// LowerHighPriorityThreadPoolCPUPriority lower CPU priority for high priority
// thread pool.
func (env *Env) LowerHighPriorityThreadPoolCPUPriority() {
	C.rocksdb_env_lower_high_priority_thread_pool_cpu_priority(env.c)
}

// Destroy deallocates the Env object.
func (env *Env) Destroy() {
	C.rocksdb_env_destroy(env.c)
	env.c = nil
}
