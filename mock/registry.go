package mock

import (
	"strings"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewRegistry returns a mock Registry.
func NewRegistry(b brazier.Backend) *Registry {
	return &Registry{
		Buckets: make(map[string]brazier.BucketConfig),
		Backend: b,
	}
}

// Registry is a mock Registry.
type Registry struct {
	Buckets           map[string]brazier.BucketConfig
	Backend           brazier.Backend
	index             []string
	CreateInvoked     bool
	BucketInvoked     bool
	BucketInfoInvoked bool
	CloseInvoked      bool
	ListInvoked       bool
}

// Create a bucket.
func (r *Registry) Create(path ...string) error {
	r.CreateInvoked = true
	name := strings.Join(path, "/")

	if _, ok := r.Buckets[name]; ok {
		return store.ErrAlreadyExists
	}

	r.Buckets[name] = brazier.BucketConfig{
		Path: path,
	}
	r.index = append(r.index, name)
	return nil
}

// BucketConfig returns the bucket informations associated with the given name.
func (r *Registry) BucketConfig(path ...string) (*brazier.BucketConfig, error) {
	r.BucketInfoInvoked = true
	name := strings.Join(path, "/")
	b, ok := r.Buckets[name]
	if !ok {
		return nil, store.ErrNotFound
	}

	return &b, nil
}

// Bucket returns the bucket associated with the given name.
func (r *Registry) Bucket(path ...string) (brazier.Bucket, error) {
	r.BucketInvoked = true
	name := strings.Join(path, "/")
	_, ok := r.Buckets[name]
	if !ok {
		return nil, store.ErrNotFound
	}

	return r.Backend.Bucket(name)
}

// List buckets.
func (r *Registry) List() ([]string, error) {
	r.ListInvoked = true
	return r.index, nil
}

// Close the Registry.
func (r *Registry) Close() error {
	r.CloseInvoked = true
	return nil
}
