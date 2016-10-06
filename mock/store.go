package mock

import "github.com/asdine/brazier"

// NewStore returns a BoltDB store
func NewStore() *Store {
	return &Store{
		Buckets: make(map[string]brazier.Bucket),
	}
}

// Store is a BoltDB store
type Store struct {
	Buckets       map[string]brazier.Bucket
	BucketInvoked bool
	CloseInvoked  bool
}

// Bucket returns the bucket associated with the given name
func (s *Store) Bucket(name string) (brazier.Bucket, error) {
	s.BucketInvoked = true
	b, ok := s.Buckets[name]
	if !ok {
		s.Buckets[name] = NewBucket()
		b = s.Buckets[name]
	}

	return b, nil
}

// Close the store
func (s *Store) Close() error {
	s.CloseInvoked = true
	return nil
}
