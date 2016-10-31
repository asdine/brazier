package mock

import (
	"strings"

	"github.com/asdine/brazier"
)

// NewBackend returns a mock backend.
func NewBackend() *Backend {
	return &Backend{
		Buckets: make(map[string]brazier.Bucket),
	}
}

// Backend is a mock backend.
type Backend struct {
	Buckets       map[string]brazier.Bucket
	BucketInvoked bool
	CloseInvoked  bool
}

// Bucket returns the bucket associated with the given name.
func (s *Backend) Bucket(path ...string) (brazier.Bucket, error) {
	s.BucketInvoked = true
	name := strings.Join(path, "/")
	b, ok := s.Buckets[name]
	if !ok {
		s.Buckets[name] = NewBucket()
		b = s.Buckets[name]
	}

	return b, nil
}

// Close the backend.
func (s *Backend) Close() error {
	s.CloseInvoked = true
	return nil
}
