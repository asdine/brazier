package mock

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewRegistry returns a mock Registry.
func NewRegistry(b brazier.Backend) *Registry {
	return &Registry{
		Backend: b,
	}
}

type bucketMeta struct {
	name     string
	children []*bucketMeta
}

// Registry is a mock Registry.
type Registry struct {
	Buckets       []*bucketMeta
	Backend       brazier.Backend
	index         []string
	CreateInvoked bool
	BucketInvoked bool
	CloseInvoked  bool
}

// Create a bucket.
func (r *Registry) Create(nodes ...string) error {
	r.CreateInvoked = true

	buckets := &r.Buckets
	var found bool

	for _, node := range nodes {
		found = false

		for _, b := range *buckets {
			if b.name == node {
				buckets = &b.children
				found = true
				break
			}
		}

		if !found {
			b := &bucketMeta{
				name: node,
			}
			*buckets = append(*buckets, b)
			buckets = &b.children
		}
	}

	if found {
		return store.ErrAlreadyExists
	}

	return nil
}

// Bucket returns the bucket associated with the given path.
func (r *Registry) Bucket(nodes ...string) (brazier.Bucket, error) {
	r.BucketInvoked = true

	buckets := r.Buckets
	var found *bucketMeta

	for _, node := range nodes {
		found = nil

		for _, b := range buckets {
			if b.name == node {
				found = b
				buckets = b.children
				break
			}
		}
	}

	if found == nil {
		return nil, store.ErrNotFound
	}

	return r.Backend.Bucket(nodes...)
}

// Close the Registry.
func (r *Registry) Close() error {
	r.CloseInvoked = true
	return nil
}
