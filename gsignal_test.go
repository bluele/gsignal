package gsignal_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/bluele/gsignal"
)

func TestWatchSignal(t *testing.T) {
	watcher := gsignal.NewWatcher()
	receive := make(chan bool, 1)
	watcher.
		Watch(func(sig os.Signal) {
			receive <- true
		}, syscall.SIGALRM).
		Watch(func(sig os.Signal) {
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
