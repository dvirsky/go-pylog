// A simple logging module that mimics the behavior of Python's logging module.
//
// All it does basically is wrap Go's logger with nice multi-level logging calls, and
// allows you to set the logging level of your app in runtime.
//
// Logging is done just like calling fmt.Sprintf:
// 		logging.Info("This object is %s and that is %s", obj, that)
//
// example output:
//	2013/05/07 01:20:26 INFO @ db.go:528: Registering plugin REPLICATION
//	2013/05/07 01:20:26 INFO @ db.go:562: Registered 6 plugins and 22 commands
//	2013/05/07 01:20:26 INFO @ slave.go:277: Running replication watchdog loop!
//	2013/05/07 01:20:26 INFO @ redis.go:49: Redis adapter listening on 0.0.0.0:2000
//	2013/05/07 01:20:26 WARN @ main.go:69: Starting adapter...
//	2013/05/07 01:20:26 INFO @ db.go:966: Finished dump load. Loaded 2 objects from dump
//	2013/05/07 01:22:26 INFO @ db.go:329: Checking persistence... 0 changes since 2m0.000297531s
//	2013/05/07 01:22:26 INFO @ db.go:337: No need to save the db. no changes...
//	2013/05/07 01:22:26 DEBUG @ db.go:341: Sleeping for 2m0s
//
package logging

import (
	"log"
	"fmt"
	"io"
	"runtime/debug"
	"path"
	"runtime"
)



const (
	DEBUG = 1
	INFO = 2
	WARNING = 4
	WARN = 4
	ERROR = 8
	CRITICAL  = 16
	QUIET = ERROR | CRITICAL  //setting for errors only
	NORMAL = INFO | WARN | ERROR | CRITICAL // default setting - all besides debug
	ALL = 255
	NOTHING = 0
)

//default logging level is ALL
var level int = ALL

// Set the logging level.
//
// Contrary to Python that specifies a minimal level, this logger is set with a bit mask
// of active levels.
//
// e.g. for INFO and ERROR use:
// 		SetLevel(logging.INFO | logging.ERROR)
//
// For everything but debug and info use:
// 		SetLevel(logging.ALL &^ (logging.INFO | logging.DEBUG))
//
func SetLevel(l int) {
	level = l
}

// Set the output writer. for now it just wraps log.SetOutput()
func SetOutPut(w io.Writer) {
	log.SetOutput(w)
}

//get the stack (line + file) context to return the caller to the log
func getContext() (file string, line int) {

	_, file, line, _ = runtime.Caller(3)
	file = path.Base(file)

	return
}

//Output debug logging messages
func Debug(msg string, args ...interface{}) {
	if level & DEBUG != 0 {
		log.Printf(fmt.Sprintf("DEBUG: %s",  msg), args...)
	}
}

//format the message
func writeMessage(level string, msg string, args ...interface {} ) {
	f, l := getContext()
	log.Printf(fmt.Sprintf("%s @ %s:%d: %s", level, f, l, msg), args...)
}

//output INFO level messages
func Info(msg string, args ...interface{}) {

	if level & INFO != 0 {

		writeMessage("INFO", msg, args...)

	}
}

//output WARNING level messages
func Warning(msg string, args ...interface{}) {
	if level & WARN != 0 {
		writeMessage("WARNING", msg, args...)
	}
}

//output ERROR level messages
func Error(msg string, args ...interface{}) {
	if level & ERROR != 0 {
		writeMessage("ERROR", msg, args...)
	}
}

//Output a CRITICAL level message while showing a stack trace
func Critical(msg string, args ...interface{}) {
	if level & CRITICAL != 0 {
		writeMessage("CRITICAL", msg, args...)
		log.Println(string(debug.Stack()))
	}
}

// Raise a PANIC while writing the stack trace to the log
func Panic(msg string, args ...interface{}) {
	log.Println(string(debug.Stack()))
	log.Panicf(msg, args...)

}

