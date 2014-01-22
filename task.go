package task

import (
	"fmt"
	"time"
)

type Handler func(*Task) Status

type Task struct {
	// Basics
	Env     *Env
	Handler Handler
	Name    string
	Args    []string
	Status  Status

	// Timing
	Start time.Time
	Stop  time.Time

	// Streams
	InputStream  *InputStream
	OutputStream *OutputStream
	ErrorStream  *OutputStream
}

// Duration returns the time the Task took to run.
func (t *Task) Duration() time.Duration {
	if t.Start.IsZero() || t.Stop.IsZero() {
		return 0
	}

	return t.Stop.Sub(t.Start)
}

func (t *Task) HasArgs() bool {
	return len(t.Args) > 0
}

func (t *Task) IsOk() bool {
    return t.Status == StatusOk
}

func (t *Task) IsError() bool {
    return t.Status == StatusErr
}

// Run executes the Task. Start/Stop time of the process are captured
// for later inspection. Closing of OutputStreams will wait for any goroutines
// consuming those streams to exit.
func (t *Task) Run() error {
	if t.Handler == nil {
		return fmt.Errorf("command not found: %s", t.Name)
	}

	t.Start = time.Now()
	t.Status = t.Handler(t)
	t.Stop = time.Now()

	if err := t.OutputStream.Close(); err != nil {
		return err
	}

	if err := t.ErrorStream.Close(); err != nil {
		return err
	}

	return nil
}
