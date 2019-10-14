package grocksdb

// rocksdb_create_iterators
// rocksdb_writebatch_putv
// rocksdb_writebatch_putv_cf
// rocksdb_writebatch_mergev
// rocksdb_writebatch_mergev_cf
// rocksdb_writebatch_deletev
// rocksdb_writebatch_deletev_cf
// rocksdb_writebatch_delete_rangev
// rocksdb_writebatch_delete_rangev_cf
// rocksdb_writebatch_iterate
// rocksdb_writebatch_set_save_point
// rocksdb_writebatch_rollback_to_save_point
// rocksdb_writebatch_pop_save_point
// rocksdb_writebatch_wi_create
// rocksdb_writebatch_wi_create_from
// rocksdb_writebatch_wi_destroy
// rocksdb_writebatch_wi_clear
// rocksdb_writebatch_wi_count
// rocksdb_writebatch_wi_put
// rocksdb_writebatch_wi_put_cf
// rocksdb_writebatch_wi_putv
// rocksdb_writebatch_wi_putv_cf
// rocksdb_writebatch_wi_merge
// rocksdb_writebatch_wi_merge_cf
// rocksdb_writebatch_wi_mergev
// rocksdb_writebatch_wi_mergev_cf
// rocksdb_writebatch_wi_delete
// rocksdb_writebatch_wi_delete_cf
// rocksdb_writebatch_wi_deletev
// rocksdb_writebatch_wi_deletev_cf
// rocksdb_writebatch_wi_delete_range
// rocksdb_writebatch_wi_delete_range_cf
// rocksdb_writebatch_wi_delete_rangev
// rocksdb_writebatch_wi_delete_rangev_cf
// rocksdb_writebatch_wi_put_log_data
// rocksdb_writebatch_wi_iterate
// rocksdb_writebatch_wi_data
// rocksdb_writebatch_wi_set_save_point
// rocksdb_writebatch_wi_rollback_to_save_point
// rocksdb_writebatch_wi_get_from_batch
// rocksdb_writebatch_wi_get_from_batch_cf
// rocksdb_writebatch_wi_get_from_batch_and_db
// rocksdb_writebatch_wi_get_from_batch_and_db_cf
// rocksdb_write_writebatch_wi
// rocksdb_writebatch_wi_create_iterator_with_base
// rocksdb_writebatch_wi_create_iterator_with_base_cf

// rocksdb_set_perf_level
// rocksdb_perfcontext_create
// rocksdb_perfcontext_reset
// rocksdb_perfcontext_report
// rocksdb_perfcontext_metric
// rocksdb_perfcontext_destroy

// rocksdb_cache_set_capacity
// rocksdb_create_mem_env
// rocksdb_env_join_all_threads
// rocksdb_env_lower_thread_pool_io_priority
// rocksdb_env_lower_high_priority_thread_pool_io_priority
// rocksdb_env_lower_thread_pool_cpu_priority
// rocksdb_env_lower_high_priority_thread_pool_cpu_priority
// rocksdb_sstfilewriter_create_with_comparator
// rocksdb_sstfilewriter_put
// rocksdb_sstfilewriter_merge
// rocksdb_sstfilewriter_delete
// rocksdb_sstfilewriter_file_size
// rocksdb_try_catch_up_with_primary
// rocksdb_livefiles_entries
// rocksdb_livefiles_deletions
// rocksdb_transactiondb_create_column_family
// rocksdb_transactiondb_open_column_families
// rocksdb_transaction_set_savepoint
// rocksdb_transaction_rollback_to_savepoint
// rocksdb_transaction_get_snapshot
// rocksdb_transaction_get_cf
// rocksdb_transaction_get_for_update_cf
// rocksdb_transactiondb_get_cf
// rocksdb_transaction_put_cf
// rocksdb_transactiondb_put_cf
// rocksdb_transactiondb_write
// rocksdb_transaction_merge
// rocksdb_transaction_merge_cf
// rocksdb_transactiondb_merge
// rocksdb_transactiondb_merge_cf
// rocksdb_transaction_delete_cf
// rocksdb_transactiondb_delete_cf
// rocksdb_transaction_create_iterator_cf
// rocksdb_transactiondb_create_iterator
// rocksdb_transactiondb_create_iterator_cf
// rocksdb_optimistictransactiondb_open
// rocksdb_optimistictransactiondb_open_column_families
// rocksdb_optimistictransactiondb_get_base_db
// rocksdb_optimistictransactiondb_close_base_db
// rocksdb_optimistictransaction_begin
// rocksdb_optimistictransactiondb_close
// rocksdb_transactiondb_options_create
// rocksdb_optimistictransaction_options_create
// rocksdb_optimistictransaction_options_destroy
// rocksdb_optimistictransaction_options_set_set_snapshot
// rocksdb_get_pinned_cf
