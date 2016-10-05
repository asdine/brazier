package mock

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewStore returns a BoltDB store
func NewStore() *Store {
	return &Store{
		Buckets: make(map[string]brazier.Bucket),
	}
}

// Store is a BoltDB store
type Store struct {
	Buckets       map[string]brazier.Bucket
	CreateInvoked bool
	BucketInvoked bool
	CloseInvoked  bool
}

// Create a bucket
func (s *Store) Create(name string) error {
	s.CreateInvoked = true
	s.Buckets[name] = NewBucket()
	return nil
}

// Bucket returns the bucket associated with the given name
func (s *Store) Bucket(name string) (brazier.Bucket, error) {
	s.BucketInvoked = true
	b, ok := s.Buckets[name]
	if !ok {
		return nil, store.ErrNotFound
	}

	return b, nil
}

// Close the store
func (s *Store) Close() error {
	s.CloseInvoked = true
	return nil
}
