package system

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestNewSignalListenerReturnsAUsableChannel(t *testing.T) {
	l := NewSignalListener()
	defer l.Stop()

	ch := l.Listen()
	if ch == nil {
		t.Fatal("Listen returned a nil channel")
	}

	// no signal has been delivered, so the channel is empty
	select {
	case s, ok := <-ch:
		t.Fatalf("received %v (open = %t) before any signal", s, ok)
	default:
	}
}

func TestListenReturnsTheSameChannel(t *testing.T) {
	l := NewSignalListener(syscall.SIGTERM)
	defer l.Stop()

	a, b := l.Listen(), l.Listen()
	if a != b {
		t.Fatal("Listen returned two different channels")
	}
}

func TestStopClosesTheChannel(t *testing.T) {
	l := NewSignalListener(syscall.SIGTERM)
	ch := l.Listen()

	l.Stop()

	select {
	case s, ok := <-ch:
		if ok {
			t.Fatalf("received %v, want a closed channel", s)
		}
	case <-time.After(time.Second):
		t.Fatal("the channel was not closed by Stop")
	}

	// the channel handed out after Stop is the closed one
	if _, ok := <-l.Listen(); ok {
		t.Fatal("Listen returned an open channel after Stop")
	}
}

func TestStopIsIdempotent(t *testing.T) {
	l := NewSignalListener(syscall.SIGTERM)

	for range 3 {
		l.Stop()
	}
}

func TestConcurrentStop(t *testing.T) {
	l := NewSignalListener(syscall.SIGTERM)

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Stop()
		}()
	}
	wg.Wait()

	if _, ok := <-l.Listen(); ok {
		t.Fatal("the channel is still open after concurrent stops")
	}
}

func TestListenersAreIndependent(t *testing.T) {
	first := NewSignalListener(syscall.SIGTERM)
	second := NewSignalListener(syscall.SIGTERM)
	defer second.Stop()

	if first.Listen() == second.Listen() {
		t.Fatal("two listeners share a channel")
	}

	first.Stop()

	if _, ok := <-first.Listen(); ok {
		t.Fatal("the stopped listener is still open")
	}
	select {
	case s, ok := <-second.Listen():
		t.Fatalf("the second listener received %v (open = %t)", s, ok)
	default:
	}
}

// guard against the zero value of os.Signal leaking through the constructor
func TestNewSignalListenerWithExplicitSignals(t *testing.T) {
	var signals []os.Signal = []os.Signal{syscall.SIGTERM, syscall.SIGINT}

	l := NewSignalListener(signals...)
	defer l.Stop()

	if l.Listen() == nil {
		t.Fatal("Listen returned a nil channel")
	}
}
