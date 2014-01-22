package task

import (
	"io"
	"sync"
)

type InputStream struct {
	Mutex  sync.Mutex
	Reader io.Reader
}

// Set assigns an io.Reader to read from as input.
func (s *InputStream) Set(r io.Reader) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Reader = r
}

// Read reads from some byte data.
func (s *InputStream) Read(data []byte) (int, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.Reader == nil {
		return 0, io.EOF
	}

	return s.Reader.Read(data)
}

// CopyTo funnels the stream into some destination io.Writer.
func (s *InputStream) CopyTo(dest io.Writer) (int64, error) {
	return io.Copy(dest, s)
}

// Close closes the stream if its able to be closed.
func (s *InputStream) Close() error {
	if closer, ok := s.Reader.(io.ReadCloser); ok {
		return closer.Close()
	}
	return nil
}
