package fs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/brazier"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

const name = "fs"

// NewStore returns a file system store
func NewStore(path string) *Store {
	return &Store{
		root: path,
	}
}

// Store is a file system store
type Store struct {
	root string
}

// Name of the store
func (s *Store) Name() string {
	return name
}

// Create a bucket and return its informations
func (s *Store) Create(id string) (*brazier.BucketInfo, error) {
	if id == "" {
		id = uniuri.NewLen(10)
	}

	err := os.Mkdir(filepath.Join(s.root, id), 0700)
	if err != nil {
		return nil, errors.Wrap(err, "fsStore.Create failed to create directory")
	}

	return &brazier.BucketInfo{
		ID:        id,
		Store:     s.Name(),
		CreatedAt: time.Now(),
	}, nil
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	return newBucket(filepath.Join(s.root, id)), nil
}
