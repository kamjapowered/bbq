package system

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// listen for os signals and expose them through a channel
//
// this type is not safe for concurrent construction, but safe for concurrent use
type SignalListener struct {
	ch   chan os.Signal
	once sync.Once
}

// create a signal listener for the given signals
//
// if no signals are provided, it listens for sigint and sigterm
//
// time: O(1)
func NewSignalListener(signals ...os.Signal) *SignalListener {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	l := &SignalListener{
		ch: make(chan os.Signal, 1),
	}

	signal.Notify(l.ch, signals...)
	return l
}

// return the channel that receives os signals
//
// receiving from this channel blocks until a signal is delivered or the listener is stopped
//
// time: O(1)
func (l *SignalListener) Listen() <-chan os.Signal {
	return l.ch
}

// stop listening for signals and release resources
//
// calling stop multiple times is safe
//
// after stop, the channel is closed
//
// time: O(1)
func (l *SignalListener) Stop() {
	l.once.Do(func() {
		signal.Stop(l.ch)
		close(l.ch)
	})
}
