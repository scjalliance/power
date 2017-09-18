package main

import (
	"os"
	"os/signal"
	"syscall"
)

// Shutdown is closed when a shutdown has been requested by the operating
// system.
type Shutdown <-chan struct{}

// NewShutdown sets up a shutdown channel that will be signaled when SIGINT
// or SIGTERM has been sent by the operating system.
func NewShutdown() (shutdown Shutdown) {
	s := make(chan struct{})
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		close(s)
	}()
	return s
}

// Signaled returns true if the shutdown has been signaled.
func (s Shutdown) Signaled() bool {
	select {
	case <-s:
		return true
	default:
		return false
	}
}
