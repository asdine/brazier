package store

import (
	"fmt"
	"path"
	"strings"

	"github.com/asdine/brazier"
)

// GetBucketOrCreate returns an existing bucket or creates it if it doesn't exist.
func GetBucketOrCreate(r brazier.Registry, path ...string) (brazier.Bucket, error) {
	if len(path) == 0 {
		return nil, ErrForbidden
	}

	bucket, err := r.Bucket(path...)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		err = r.Create(path...)
		if err != nil {
			return nil, err
		}
		bucket, err = r.Bucket(path...)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}

// SplitPathKey returns a slice of bucket names, representing the path of a bucket,
// and a key.
func SplitPathKey(rawPath string) ([]string, string) {
	nodes := splitPath(rawPath)
	if len(nodes) == 0 {
		return nodes, ""
	}

	return nodes[:len(nodes)-1], nodes[len(nodes)-1]
}

func splitPath(rawPath string) []string {
	rawPath = path.Clean(path.Join("/", rawPath))
	nodes := strings.Split(rawPath, "/")

	var list []string
	for _, n := range nodes {
		if len(n) == 0 {
			continue
		}

		list = append(list, n)
	}
	return list
}

// NewStore instantiates a new Store with the given Registry.
func NewStore(r brazier.Registry) *Store {
	return &Store{
		r: r,
	}
}

// A Store manages items from various backends.
type Store struct {
	r brazier.Registry
}

// CreateBucket creates a bucket at the given path.
func (s *Store) CreateBucket(rawPath string) error {
	nodes := splitPath(rawPath)
	if len(nodes) == 0 {
		return ErrAlreadyExists
	}

	fmt.Printf("`%s`: `%v`", rawPath, nodes)
	return s.r.Create(nodes...)
}

// Save the value at the given path.
func (s *Store) Save(rawPath string, value []byte) (*brazier.Item, error) {
	path, key := SplitPathKey(rawPath)
	bucket, err := GetBucketOrCreate(s.r, path...)
	if err != nil {
		return nil, err
	}

	i, err := bucket.Save(key, value)
	bucket.Close()
	return i, err
}

// Get returns the item saved at the given path.
func (s *Store) Get(rawPath string) (*brazier.Item, error) {
	path, key := SplitPathKey(rawPath)
	bucket, err := s.r.Bucket(path...)
	if err != nil {
		return nil, err
	}

	i, err := bucket.Get(key)
	bucket.Close()
	return i, err
}

// List the content of the bucket.
func (s *Store) List(rawPath string, page int, perPage int) ([]brazier.Item, error) {
	nodes := splitPath(rawPath)
	bucket, err := s.r.Bucket(nodes...)
	if err != nil {
		return nil, err
	}

	list, err := bucket.Page(page, perPage)
	bucket.Close()
	return list, err
}

// Delete the key from the bucket.
func (s *Store) Delete(rawPath string) error {
	path, key := SplitPathKey(rawPath)
	bucket, err := s.r.Bucket(path...)
	if err != nil {
		return err
	}

	err = bucket.Delete(key)
	bucket.Close()
	return err
}
