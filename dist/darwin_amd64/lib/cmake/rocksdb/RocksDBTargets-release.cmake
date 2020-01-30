#----------------------------------------------------------------
# Generated CMake target import file for configuration "Release".
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "RocksDB::rocksdb" for configuration "Release"
set_property(TARGET RocksDB::rocksdb APPEND PROPERTY IMPORTED_CONFIGURATIONS RELEASE)
set_target_properties(RocksDB::rocksdb PROPERTIES
  IMPORTED_LINK_INTERFACE_LANGUAGES_RELEASE "CXX"
  IMPORTED_LOCATION_RELEASE "${_IMPORT_PREFIX}/lib/librocksdb.a"
  )

list(APPEND _IMPORT_CHECK_TARGETS RocksDB::rocksdb )
list(APPEND _IMPORT_CHECK_FILES_FOR_RocksDB::rocksdb "${_IMPORT_PREFIX}/lib/librocksdb.a" )

# Import target "RocksDB::rocksdb-shared" for configuration "Release"
set_property(TARGET RocksDB::rocksdb-shared APPEND PROPERTY IMPORTED_CONFIGURATIONS RELEASE)
set_target_properties(RocksDB::rocksdb-shared PROPERTIES
  IMPORTED_LOCATION_RELEASE "${_IMPORT_PREFIX}/lib/librocksdb.6.6.3.dylib"
  IMPORTED_SONAME_RELEASE "@rpath/librocksdb.6.dylib"
  )

list(APPEND _IMPORT_CHECK_TARGETS RocksDB::rocksdb-shared )
list(APPEND _IMPORT_CHECK_FILES_FOR_RocksDB::rocksdb-shared "${_IMPORT_PREFIX}/lib/librocksdb.6.6.3.dylib" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)
