PyLog
====


A simple logging module that mimics the behavior of Python's logging module.

All it does basically is wrap Go's logger with nice multi-level logging calls, and
allows you to set the logging level of your app in runtime.

Logging is done just like calling fmt.Sprintf:

```go
logging.Info("This object is %s and that is %s", obj, that)
```

Logging level can be set to whatever you want it to be, in runtime. Contrary to Python that specifies a minimal level, this logger is set with a bit mask of active levels.

```go
//for INFO and ERROR use:
logging.SetLevel(logging.INFO | logging.ERROR)

// For everything but debug and info use:
logging.SetLevel(logging.ALL &^ (logging.INFO | logging.DEBUG))
```

### Installation:

```
go get github.com/dvirsky/go-pylog/logging
```

### Usage Example:

```go
package main

import (
	"github.com/dvirsky/go-pylog/logging"
)

func main() {

	logging.Info("All Your Base Are Belong to %s!", "us")

	logging.Critical("And now with a stack trace")
}
```




### Example Output:

```
2013/05/07 01:20:26 INFO @ db.go:528: Registering plugin REPLICATION
2013/05/07 01:20:26 INFO @ db.go:562: Registered 6 plugins and 22 commands
2013/05/07 01:20:26 INFO @ slave.go:277: Running replication watchdog loop!
2013/05/07 01:20:26 INFO @ redis.go:49: Redis adapter listening on 0.0.0.0:2000
2013/05/07 01:20:26 WARN @ main.go:69: Starting adapter...
2013/05/07 01:20:26 INFO @ db.go:966: Finished dump load. Loaded 2 objects from dump
2013/05/07 01:22:26 INFO @ db.go:329: Checking persistence... 0 changes since 2m0.000297531s
2013/05/07 01:22:26 DEBUG @ db.go:341: Sleeping for 2m0s
```
