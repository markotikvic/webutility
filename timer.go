package webutility

import (
	"fmt"
	"time"
)

// Timer ...
type Timer struct {
	name         string
	running      bool
	started      time.Time
	stopped      time.Time
	lastDuration time.Duration
}

// NewTimer ...
func NewTimer(name string) *Timer {
	t := &Timer{name: name}
	t.Reset()
	return t
}

// Start ...
func (t *Timer) Start() time.Time {
	t.running = true
	t.started = time.Now()
	return t.started
}

// Stop ...
func (t *Timer) Stop() time.Duration {
	t.running = false
	t.stopped = time.Now()
	t.lastDuration = t.stopped.Sub(t.started)
	return t.lastDuration
}

// LastRunDuration ...
func (t *Timer) LastRunDuration() time.Duration {
	return t.lastDuration
}

// Clear ...
func (t *Timer) Clear() {
	t.started = time.Now()
	t.stopped = time.Now()
	t.lastDuration = 0
}

// Reset ...
func (t *Timer) Reset() {
	t.Stop()
	t.Start()
}

// Elapsed ...
func (t *Timer) Elapsed() time.Duration {
	if t.running {
		return time.Now().Sub(t.started)
	}
	return 0
}

// Print ...
func (t *Timer) Print(s string) {
	status := "RUNNING"
	if !t.running {
		status = "STOPPED"
	}
	fmt.Printf("timer[%s][%s]: %v %s\n", t.name, status, t.Elapsed(), s)
}

// Lap ...
func (t *Timer) Lap(s string) {
	t.Print(s)
	t.Reset()
}
