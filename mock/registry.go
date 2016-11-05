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
	Buckets         []*bucketMeta
	Backend         brazier.Backend
	index           []string
	CreateInvoked   bool
	BucketInvoked   bool
	CloseInvoked    bool
	ChildrenInvoked bool
}

// Create a bucket.
func (r *Registry) Create(nodes ...string) error {
	r.CreateInvoked = true

	if len(nodes) == 0 {
		return store.ErrForbidden
	}

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

	if len(nodes) == 0 {
		return nil, store.ErrForbidden
	}

	_, err := r.bucket(nodes...)
	if err != nil {
		return nil, err
	}

	return r.Backend.Bucket(nodes...)
}

func (r *Registry) bucket(nodes ...string) (*bucketMeta, error) {
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

	return found, nil
}

// Children buckets of the specified path.
func (r *Registry) Children(nodes ...string) ([]brazier.Item, error) {
	b, err := r.bucket(nodes...)
	if err != nil {
		return nil, err
	}

	return r.children(b)
}

func (r *Registry) children(b *bucketMeta) ([]brazier.Item, error) {
	var tree []brazier.Item

	for _, child := range b.children {
		items, err := r.children(child)
		if err != nil {
			return nil, err
		}

		tree = append(tree, brazier.Item{
			Key:      child.name,
			Children: items,
		})
	}

	return tree, nil
}

// Close the Registry.
func (r *Registry) Close() error {
	r.CloseInvoked = true
	return nil
}
