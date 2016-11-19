package mock

import "github.com/asdine/brazier"

// NewBackend returns a mock backend.
func NewBackend() *Backend {
	return &Backend{
		Tree: NewBucket("/"),
	}
}

// Backend is a mock backend.
type Backend struct {
	Tree          *Bucket
	BucketInvoked bool
	CloseInvoked  bool
}

// Bucket returns the bucket associated with the given path.
func (s *Backend) Bucket(nodes ...string) (brazier.Bucket, error) {
	s.BucketInvoked = true

	if len(nodes) == 0 {
		return s.Tree, nil
	}

	var b *Bucket
	buckets := &s.Tree.Children

	var found bool
	for _, node := range nodes {
		found = false

		for _, b = range *buckets {
			if b.Name == node {
				buckets = &b.Children
				found = true
				break
			}
		}

		if !found {
			b = NewBucket(node)
			*buckets = append(*buckets, b)
			buckets = &b.Children
		}
	}

	return b, nil
}

// Close the backend.
func (s *Backend) Close() error {
	s.CloseInvoked = true
	return nil
}
