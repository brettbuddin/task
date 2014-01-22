package task

import (
	"io"
	"sync"
)

type OutputStream struct {
	Mutex   sync.Mutex
	Writers []io.Writer
	Waiting sync.WaitGroup
}

// Tail fans-in all writing streams to a single output destination.
func (s *OutputStream) Consume(dest io.Writer) {
	s.Waiting.Add(1)
	out, in := io.Pipe()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Writers = append(s.Writers, in)

	go func() {
		io.Copy(dest, out)
		s.Waiting.Done()
	}()
}

// Write writes some byte data to all tailing goroutines.
func (s *OutputStream) Write(data []byte) (int, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, writer := range s.Writers {
		if n, err := writer.Write(data); err != nil {
			return n, err
		}
	}

	return len(data), nil
}

// Close closes all writer streams if they are able to be closed.
// Waits for all goroutines tailing the stream to exit.
func (s *OutputStream) Close() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var first error
	for _, writer := range s.Writers {
		if closer, ok := writer.(io.WriteCloser); ok {
			err := closer.Close()

			if err != nil && first != nil {
				first = err
			}
		}
	}

	s.Waiting.Wait()
	return first
}
