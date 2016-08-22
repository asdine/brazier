package fs

import (
	"os"
	"path/filepath"

	"github.com/asdine/brazier"
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
func (s *Store) Create(id string) error {
	err := os.Mkdir(filepath.Join(s.root, id), 0700)
	return errors.Wrap(err, "fsStore.Create failed to create directory")
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	return newBucket(filepath.Join(s.root, id)), nil
}
