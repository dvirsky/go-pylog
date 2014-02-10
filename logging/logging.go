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
	"fmt"
	//	"github.com/samuel/go-thrift/examples/scribe"
	//"github.com/samuel/go-thrift/thrift"
	"io"
	"log"
	"strings"
	//	"net"
	"path"
	"runtime"
	"runtime/debug"
)

const (
	DEBUG    = 1
	INFO     = 2
	WARNING  = 4
	WARN     = 4
	ERROR    = 8
	CRITICAL = 16
	QUIET    = ERROR | CRITICAL               //setting for errors only
	NORMAL   = INFO | WARN | ERROR | CRITICAL // default setting - all besides debug
	ALL      = 255
	NOTHING  = 0
)

var levels_ascending = []int{DEBUG, INFO, WARNING, ERROR, CRITICAL}

var LevlelsByName = map[string]int{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"WARNING":  WARN,
	"WARN":     WARN,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
	"QUIET":    QUIET,
	"NORMAL":   NORMAL,
	"ALL":      ALL,
	"NOTHING":  NOTHING,
}

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

// Set a minimal level for loggin, setting all levels higher than this level as well.
//
// the severity order is DEBUG, INFO, WARNING, ERROR, CRITICAL
func SetMinimalLevel(l int) {

	newLevel := 0
	for _, level := range levels_ascending {
		if level >= l {

			newLevel |= level
			fmt.Println(level, newLevel)
		}
	}
	SetLevel(newLevel)

}

// Set minimal level by string, useful for config files and command line arguments. Case insensitive.
//
// Possible level names are DEBUG, INFO, WARNING, ERROR, CRITICAL
func SetMinimalLevelByName(l string) error {
	l = strings.ToUpper(strings.Trim(l, " "))
	level, found := LevlelsByName[l]
	if !found {
		Error("Could not set level - not found level %s", l)
		return fmt.Errorf("Invalid level %s", l)
	}

	SetMinimalLevel(level)
	return nil
}

// Set the output writer. for now it just wraps log.SetOutput()
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

//a pluggable logger interface
type LoggingHandler interface {
	Emit(level, file string, line int, message string, args ...interface{}) error
}

var formatString = "%[1]s @ %[2]s:%[3]d: %[4]s"

// Set the logger's format string. The arguments passed to it are always "level, file string, line int, message string"
//
// This means that if you want to change the order they appear, you should use explicit index numbers in formatting.
//
// The default format is "%[1]s @ %[2]s:%[2]d: %[4]s"
func SetFormatString(format string) {
	formatString = format
}

func GetFormatString() string {
	return formatString
}

type strandardHandler struct{}

// default handling interface - just
func (l strandardHandler) Emit(level, file string, line int, message string, args ...interface{}) error {
	log.Printf(fmt.Sprintf(formatString, level, file, line, message), args...)
	return nil
}

var currentHandler LoggingHandler = strandardHandler{}

// Set the current handler of the library. We currently support one handler, but it might be nice to have more
func SetHandler(h LoggingHandler) {
	currentHandler = h
}

//get the stack (line + file) context to return the caller to the log
func getContext() (file string, line int) {

	_, file, line, _ = runtime.Caller(3)
	file = path.Base(file)

	return
}

//Output debug logging messages
func Debug(msg string, args ...interface{}) {
	if level&DEBUG != 0 {
		writeMessage("DEBUG", msg, args...)
	}
}

//format the message
func writeMessage(level string, msg string, args ...interface{}) {
	f, l := getContext()

	// We go over the args, and replace any function pointer with the signature
	// func() interface{} with the return value of executing it now.
	// This allows lazy evaluation of arguments which are return values
	for i, arg := range args {
		switch arg.(type) {
		case func() interface{}:
			args[i] = arg.(func() interface{})()
		default:

		}
	}
	err := currentHandler.Emit(level, f, l, msg, args...)
	if err != nil {
		log.Printf("Error writing log message: %s\n", err)
		log.Printf(fmt.Sprintf(formatString, level, f, l, msg), args...)
	}

}

//output INFO level messages
func Info(msg string, args ...interface{}) {

	if level&INFO != 0 {

		writeMessage("INFO", msg, args...)

	}
}

// Output WARNING level messages
func Warning(msg string, args ...interface{}) {
	if level&WARN != 0 {
		writeMessage("WARNING", msg, args...)
	}
}

// Same as Warning() but return a formatted error object, regardless of logging level
func Warningf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args)
	if level&WARN != 0 {
		writeMessage("WARNING", err.Error())
	}

	return err
}

// Output ERROR level messages
func Error(msg string, args ...interface{}) {
	if level&ERROR != 0 {
		writeMessage("ERROR", msg, args...)
	}
}

// Same as Error() but also returns a new formatted error object with the message regardless of logging level
func Errorf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args)
	if level&ERROR != 0 {
		writeMessage("ERROR", err.Error())
	}
	return err
}

// Output a CRITICAL level message while showing a stack trace
func Critical(msg string, args ...interface{}) {
	if level&CRITICAL != 0 {
		writeMessage("CRITICAL", msg, args...)
		log.Println(string(debug.Stack()))
	}
}

// Same as critical but also returns an error object with the message regardless of logging level
func Criticalf(msg string, args ...interface{}) error {

	err := fmt.Errorf(msg, args)
	if level&CRITICAL != 0 {
		writeMessage("CRITICAL", err.Error())
		log.Println(string(debug.Stack()))
	}
	return err
}

// Raise a PANIC while writing the stack trace to the log
func Panic(msg string, args ...interface{}) {
	log.Println(string(debug.Stack()))
	log.Panicf(msg, args...)

}
