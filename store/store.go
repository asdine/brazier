package store

import (
	"path"
	"strings"

	"github.com/asdine/brazier"
)

// NewStore instantiates a new Store with the given Registry.
func NewStore(r brazier.Registry) *Store {
	return &Store{
		Registry: r,
	}
}

// A Store manages items from various backends.
type Store struct {
	Registry brazier.Registry
}

// CreateBucket creates a bucket at the given path.
func (s *Store) CreateBucket(rawPath string) error {
	if len(rawPath) == 0 {
		return ErrForbidden
	}

	if rawPath[len(rawPath)-1] != '/' {
		return ErrForbidden
	}

	nodes := splitPath(rawPath)
	if len(nodes) == 0 {
		return ErrAlreadyExists
	}

	return s.Registry.Create(nodes...)
}

// Save the value at the given path.
func (s *Store) Put(rawPath string, value []byte) (*brazier.Item, error) {
	nodes, key := SplitPathKey(rawPath)
	if key == "" {
		return nil, ErrForbidden
	}

	bucket, err := GetBucketOrCreate(s.Registry, nodes...)
	if err != nil {
		return nil, err
	}

	i, err := bucket.Save(key, value)
	bucket.Close()
	return i, err
}

// Get returns the item saved at the given path.
func (s *Store) Get(rawPath string) (*brazier.Item, error) {
	nodes, key := SplitPathKey(rawPath)
	if key == "" {
		return nil, ErrForbidden
	}

	bucket, err := s.Registry.Bucket(nodes...)
	if err != nil {
		return nil, err
	}

	i, err := bucket.Get(key)
	bucket.Close()
	return i, err
}

// Delete the key from the bucket.
func (s *Store) Delete(rawPath string) error {
	nodes, key := SplitPathKey(rawPath)
	if key == "" {
		return ErrForbidden
	}

	bucket, err := s.Registry.Bucket(nodes...)
	if err != nil {
		return err
	}

	err = bucket.Delete(key)
	bucket.Close()
	return err
}

// List the content of the bucket.
func (s *Store) List(rawPath string, page int, perPage int) ([]brazier.Item, error) {
	nodes, key := SplitPathKey(rawPath)
	if key != "" {
		return nil, ErrForbidden
	}

	bucket, err := s.Registry.Bucket(nodes...)
	if err != nil {
		return nil, err
	}

	list, err := bucket.Page(page, perPage)
	bucket.Close()
	return list, err
}

// Tree returns the content of the bucket and of all its children.
func (s *Store) Tree(rawPath string) ([]brazier.Item, error) {
	nodes, key := SplitPathKey(rawPath)
	if key != "" {
		return nil, ErrForbidden
	}

	buckets, err := s.Registry.Children(nodes...)
	if err != nil {
		return nil, err
	}

	return s.tree(buckets, nodes...)
}

func (s *Store) tree(buckets []brazier.Item, nodes ...string) ([]brazier.Item, error) {
	var items []brazier.Item
	var err error

	bucket, err := s.Registry.Bucket(nodes...)
	if err != nil {
		return nil, err
	}

	items, err = bucket.Page(1, -1)
	bucket.Close()
	if err != nil {
		return nil, err
	}

	for _, b := range buckets {
		i := brazier.Item{
			Key: b.Key,
		}

		i.Children, err = s.tree(b.Children, append(nodes, b.Key)...)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// Close the registry and the backend connection.
func (s *Store) Close() error {
	return s.Registry.Close()
}

// GetBucketOrCreate returns an existing bucket or creates it if it doesn't exist.
func GetBucketOrCreate(r brazier.Registry, nodes ...string) (brazier.Bucket, error) {
	bucket, err := r.Bucket(nodes...)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		err = r.Create(nodes...)
		if err != nil {
			return nil, err
		}
		bucket, err = r.Bucket(nodes...)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}

// SplitPathKey returns a slice of bucket names, representing the path of a bucket,
// and a key.
func SplitPathKey(rawPath string) ([]string, string) {
	if rawPath == "" {
		return nil, ""
	}

	dir, key := path.Split(rawPath)
	return splitPath(dir), key
}

func splitPath(rawPath string) []string {
	if rawPath == "" {
		return nil
	}

	rawPath = strings.TrimPrefix(path.Clean(rawPath), "/")

	if rawPath == "" {
		return nil
	}

	return strings.Split(rawPath, "/")
}
