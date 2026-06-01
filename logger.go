package grocksdb

// #include "rocksdb/c.h"
// #include "grocksdb.h"
import "C"
import "unsafe"

// Logger struct.
type Logger struct {
	c *C.rocksdb_logger_t
}

func NewStderrLogger(level InfoLogLevel, prefix string) *Logger {
	prefix_ := C.CString(prefix)
	defer C.free(unsafe.Pointer(prefix_))

	return &Logger{
		c: C.rocksdb_logger_create_stderr_logger(C.int(level), prefix_),
	}
}

// LogHandler is invoked for each log message emitted by RocksDB when a
// callback logger is installed. level is the severity of the message.
type LogHandler = func(level InfoLogLevel, msg string)

// NewCallbackLogger creates a Logger that forwards every RocksDB log message
// to the provided Go handler. Messages below level are filtered out by RocksDB.
func NewCallbackLogger(level InfoLogLevel, handler LogHandler) *Logger {
	idx := registerLogHandler(handler)
	return &Logger{
		c: C.gorocksdb_logger_create_callback(C.int(level), C.uintptr_t(idx)),
	}
}

// Destroy Logger.
func (l *Logger) Destroy() {
	C.rocksdb_logger_destroy(l.c)
	l.c = nil
}

// Hold references to log handlers so they survive across the cgo boundary.
var logHandlers = NewCOWList()

func registerLogHandler(handler LogHandler) int {
	return logHandlers.Append(handler)
}

//export gorocksdb_logger_log
func gorocksdb_logger_log(idx int, level C.uint, cMsg *C.char, cLen C.size_t) {
	handler := logHandlers.Get(idx).(LogHandler)
	handler(InfoLogLevel(level), C.GoStringN(cMsg, C.int(cLen)))
}
