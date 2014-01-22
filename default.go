package task

var defaultContainer *Container

// DefaultContainer creates (if not already created) and returns
// the default package level Container.
func DefaultContainer() *Container {
	if defaultContainer == nil {
		defaultContainer = NewContainer()
	}

	return defaultContainer
}

// Environment returns the environment of the default Container.
func Environment() *Env {
	return DefaultContainer().Env
}

// Define registers a Handler with the default Container.
func Define(name string, handler Handler) {
	DefaultContainer().Define(name, handler)
}

// New creates a new Task using a Handler defined
// in the default Container.
func New(name string, args ...string) *Task {
	task := &Task{
		Handler:      DefaultContainer().Handler(name),
		Name:         name,
		Args:         args,
		InputStream:  &InputStream{},
		OutputStream: &OutputStream{},
		ErrorStream:  &OutputStream{},
	}

	return task
}
