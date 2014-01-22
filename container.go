package task

import (
	"sync"
)

// Container represents the overall environment that Tasks run within.
// Here you'll define Handlers to be "instantiated" via Tasks.
type Container struct {
	Env      *Env
	Mutex    sync.RWMutex
	Handlers map[string]Handler
}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{
		Handlers: make(map[string]Handler),
	}
}

// Define registers a new Handler that Tasks can use.
func (c *Container) Define(name string, handler Handler) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Handlers[name] = handler
}

// Handler returns a defined Handler. Nil will be returned in
// the event that the Handler has not been defined.
func (c *Container) Handler(name string) Handler {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Handlers[name]
}

// NewTask creates a new Task.
func (c *Container) NewTask(name string, args ...string) *Task {
	task := &Task{
		Env:          c.Env,
		Handler:      c.Handler(name),
		Name:         name,
		Args:         args,
		InputStream:  &InputStream{},
		OutputStream: &OutputStream{},
		ErrorStream:  &OutputStream{},
	}

	return task
}
