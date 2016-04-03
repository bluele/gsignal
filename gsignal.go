package gsignal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// type for callback handler
type CallbackHandler func(os.Signal)
type callbacks map[os.Signal]CallbackHandler

// watcher object
type Watcher struct {
	stopCh    chan byte
	callbacks callbacks
	mu        sync.RWMutex
}

// Create a new signal watcher.
func NewWatcher() *Watcher {
	gs := &Watcher{
		stopCh:    make(chan byte, 1),
		callbacks: make(callbacks),
	}
	return gs
}

// cb: callback hander when watcher catches specified signals.
// signals: signals to monitor.
func (gs *Watcher) Watch(cb CallbackHandler, signals ...os.Signal) *Watcher {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for _, sg := range signals {
		gs.callbacks[sg] = cb
	}

	return gs
}

// signals: signals to unmonitor.
func (gs *Watcher) UnWatch(signals ...os.Signal) *Watcher {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for _, sg := range signals {
		delete(gs.callbacks, sg)
	}

	return gs
}

// get callback handler of specified signal and boolean value for exists callback.
func (gs *Watcher) GetCallback(sig os.Signal) (CallbackHandler, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	cb, ok := gs.callbacks[sig]
	if ok {
		return cb, true
	}
	return nil, false
}

// start watching specified signals asynchronously.
func (gs *Watcher) Run() {
	chann := make(chan os.Signal, 1)
	gs.mu.RLock()
	sgs := make([]os.Signal, len(gs.callbacks))
	for sg := range gs.callbacks {
		sgs = append(sgs, sg)
	}
	gs.mu.RUnlock()
	signal.Notify(chann, sgs...)
	go func() {
		defer signal.Stop(chann)
		for {
			select {
			case <-gs.stopCh:
				return
			case sig := <-chann:
				gs.mu.RLock()
				cb, ok := gs.callbacks[sig]
				gs.mu.RUnlock()
				if !ok {
					continue
				}
				go cb(sig)
			}
		}
	}()
}

// stop watching signals.
func (gs *Watcher) Stop() {
	select {
	case gs.stopCh <- 0:
	default:
	}
}

// reload watching signals.
func (gs *Watcher) Reload() {
	gs.Stop()
	gs.Run()
}

// send SIGALRM to this process.
func Alarm(delay time.Duration) {
	SendSignal(delay, syscall.SIGALRM)
}

// delay: delay for sending signal.
// sigs: signal numbers
// send specified signal to this process.
func SendSignal(delay time.Duration, sigs ...syscall.Signal) {
	time.AfterFunc(delay, func() {
		pid := os.Getpid()
		for _, sig := range sigs {
			syscall.Kill(pid, sig)
		}
	})
}
