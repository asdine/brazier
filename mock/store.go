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
	NameInvoked   bool
	CreateInvoked bool
	BucketInvoked bool
}

// Name of the store
func (s *Store) Name() string {
	s.NameInvoked = true
	return "mock"
}

// Create a bucket
func (s *Store) Create(id string) error {
	s.CreateInvoked = true
	s.Buckets[id] = NewBucket()
	return nil
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	s.BucketInvoked = true
	b, ok := s.Buckets[id]
	if !ok {
		return nil, store.ErrNotFound
	}

	return b, nil
}
