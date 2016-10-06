package mock

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewRegistry returns a BoltDB Registry
func NewRegistry(s brazier.Store) *Registry {
	return &Registry{
		Buckets: make(map[string]brazier.BucketInfo),
		Store:   s,
	}
}

// Registry is a BoltDB store
type Registry struct {
	Buckets           map[string]brazier.BucketInfo
	Store             brazier.Store
	index             []string
	CreateInvoked     bool
	BucketInvoked     bool
	BucketInfoInvoked bool
	CloseInvoked      bool
	ListInvoked       bool
}

// Create a bucket
func (r *Registry) Create(name string) error {
	r.CreateInvoked = true
	r.Buckets[name] = brazier.BucketInfo{
		Name: name,
	}
	r.index = append(r.index, name)
	return nil
}

// BucketInfo returns the bucket informations associated with the given name
func (r *Registry) BucketInfo(name string) (*brazier.BucketInfo, error) {
	r.BucketInfoInvoked = true
	b, ok := r.Buckets[name]
	if !ok {
		return nil, store.ErrNotFound
	}

	return &b, nil
}

// Bucket returns the bucket associated with the given name
func (r *Registry) Bucket(name string) (brazier.Bucket, error) {
	r.BucketInvoked = true
	_, ok := r.Buckets[name]
	if !ok {
		return nil, store.ErrNotFound
	}

	return r.Store.Bucket(name)
}

// List buckets
func (r *Registry) List() ([]string, error) {
	r.ListInvoked = true
	return r.index, nil
}

// Close the store
func (r *Registry) Close() error {
	r.CloseInvoked = true
	return nil
}
