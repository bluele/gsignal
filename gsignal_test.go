package gsignal_test

import (
	"github.com/bluele/gsignal"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWatchSignal(t *testing.T) {
	watcher := gsignal.NewWatcher()
	receive := make(chan bool, 1)
	watcher.Watch(func(sig os.Signal) {
		receive <- true
	}, syscall.SIGALRM)

	watcher.Watch(func(sig os.Signal) {
		t.Error("Not expected catching signal: ", sig)
	}, syscall.SIGHUP)

	watcher.Run()
	watcher.UnWatch(syscall.SIGHUP)

	gsignal.SendSignal(time.Millisecond, syscall.SIGALRM)
	time.Sleep(time.Millisecond)

	gsignal.SendSignal(time.Millisecond, syscall.SIGHUP)
	time.Sleep(time.Millisecond)

	select {
	case <-receive:
	default:
		t.Error("channel `receive` should not be empty.")
	}

}
