package mock

import (
	"net"
	"time"

	"github.com/asdine/brazier"
)

// NewServer allocates a new Mock Server
func NewServer(r brazier.Registry, s brazier.Store) brazier.Server {
	return &Server{
		quit: make(chan struct{}),
	}
}

// Server is a mock Brazier server
type Server struct {
	ServeInvoked bool
	StopInvoked  bool
	quit         chan struct{}
}

// Serve simulates a blocking server
func (s *Server) Serve(l net.Listener) error {
	s.ServeInvoked = true
	<-s.quit
	return nil
}

// Stop stops the mock server
func (s *Server) Stop(time.Duration) {
	s.StopInvoked = true
	s.quit <- struct{}{}
}
