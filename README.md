# GSignal

Golang library for monitoring asynchronously signals.

## Getting Started

```
$ go get github.com/bluele/gsignal
```

## Examples

```go
package main

import (
  "github.com/bluele/gsignal"
  "log"
  "os"
  "syscall"
  "time"
)

func main() {
  gsg := gsignal.NewWatcher()
  gsg.Watch(func(sig os.Signal) {
    log.Println("Catch signal: ", sig)
    os.Exit(0)
  }, os.Interrupt, syscall.SIGTERM)

  gsg.Watch(func(sig os.Signal) {
    log.Println("Ignore signal: ", sig)
  }, syscall.SIGHUP)

  // start watching specified signal numbers.
  gsg.Run()

  log.Println("PID: ", os.Getpid())
  log.Println("Sleep...")
  time.Sleep(time.Minute)
}
```

### send Ctr^C to this process.

```
$ go run example.go
2015/01/04 14:30:40 PID:  96311 # your process id.
2015/01/04 14:30:40 Sleep...
^C2015/01/04 14:30:41 Catch signal:  interrupt
```

### send SIGTERM to this process.

```
$ kill -SIGTERM <your pid>
2015/01/04 14:30:40 PID:  96311 # your process id.
2015/01/04 14:30:40 Sleep...
^C2015/01/04 14:30:41 Catch signal:  terminated
```

### send SIGHUP to this process.

```
$ kill -SIGHUP <your pid>
2015/01/04 14:30:40 PID:  96311 # your process id.
2015/01/04 14:30:40 Sleep...
^C2015/01/04 14:30:41 Ignore signal:  hangup
```

## Testing

```
$ go test
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>