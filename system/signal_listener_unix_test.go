//go:build unix

package system

import (
	"os"
	"syscall"
	"testing"
	"time"
)

// send a signal to this process
//
// a keeper listener must be registered for it first, so the runtime never
// applies the default action and a missed delivery fails the test cleanly
func raise(t *testing.T, sig syscall.Signal) {
	t.Helper()

	if err := syscall.Kill(os.Getpid(), sig); err != nil {
		t.Fatalf("sending %v: %v", sig, err)
	}
}

func expect(t *testing.T, ch <-chan os.Signal, want os.Signal) {
	t.Helper()

	select {
	case got, ok := <-ch:
		if !ok {
			t.Fatalf("the channel was closed while waiting for %v", want)
		}
		if got != want {
			t.Fatalf("received %v, want %v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("%v was never delivered", want)
	}
}

func TestListenerReceivesItsSignal(t *testing.T) {
	keeper := NewSignalListener(syscall.SIGUSR1)
	defer keeper.Stop()

	l := NewSignalListener(syscall.SIGUSR1)
	defer l.Stop()

	raise(t, syscall.SIGUSR1)

	expect(t, l.Listen(), syscall.SIGUSR1)
	expect(t, keeper.Listen(), syscall.SIGUSR1)
}

func TestDefaultSignalsAreSigintAndSigterm(t *testing.T) {
	keeper := NewSignalListener(syscall.SIGTERM, syscall.SIGINT)
	defer keeper.Stop()

	l := NewSignalListener()
	defer l.Stop()

	raise(t, syscall.SIGTERM)
	expect(t, l.Listen(), syscall.SIGTERM)
	expect(t, keeper.Listen(), syscall.SIGTERM)

	raise(t, syscall.SIGINT)
	expect(t, l.Listen(), syscall.SIGINT)
	expect(t, keeper.Listen(), syscall.SIGINT)
}

func TestUnregisteredSignalsAreNotDelivered(t *testing.T) {
	keeper := NewSignalListener(syscall.SIGUSR1)
	defer keeper.Stop()

	l := NewSignalListener(syscall.SIGUSR2)
	defer l.Stop()

	raise(t, syscall.SIGUSR1)
	expect(t, keeper.Listen(), syscall.SIGUSR1)

	select {
	case s := <-l.Listen():
		t.Fatalf("received %v, want nothing", s)
	default:
	}
}

func TestStopUnsubscribesBeforeClosing(t *testing.T) {
	keeper := NewSignalListener(syscall.SIGUSR1)
	defer keeper.Stop()

	l := NewSignalListener(syscall.SIGUSR1)
	l.Stop()

	// a signal raised after Stop must not be sent to the closed channel
	raise(t, syscall.SIGUSR1)
	expect(t, keeper.Listen(), syscall.SIGUSR1)

	if _, ok := <-l.Listen(); ok {
		t.Fatal("the stopped listener received a signal")
	}
}
