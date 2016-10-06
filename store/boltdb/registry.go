package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb/internal"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// NewRegistry returns a BoltDB Registry
func NewRegistry(path string, s brazier.Store) (*Registry, error) {
	var err error

	db, err := storm.Open(
		path,
		storm.AutoIncrement(),
		storm.Codec(protobuf.Codec),
		storm.BoltOptions(0644, &bolt.Options{
			Timeout: time.Duration(50) * time.Millisecond,
		}),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open database")
	}

	return &Registry{
		DB:    db,
		Store: s,
	}, nil
}

// Registry is a BoltDB store
type Registry struct {
	DB    *storm.DB
	Store brazier.Store
}

// Create a bucket
func (r *Registry) Create(name string) error {
	err := r.DB.Save(&internal.Bucket{
		Name: name,
	})

	if err == storm.ErrAlreadyExists {
		return store.ErrAlreadyExists
	}

	return err
}

// BucketConfig returns the bucket informations associated with the given name
func (r *Registry) BucketConfig(name string) (*brazier.BucketConfig, error) {
	var b internal.Bucket
	err := r.DB.One("Name", name, &b)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.registry.Bucket failed to fetch bucket")
	}

	return &brazier.BucketConfig{
		Name: b.Name,
	}, nil
}

// Bucket returns the selected bucket from the Store
func (r *Registry) Bucket(name string) (brazier.Bucket, error) {
	info, err := r.BucketConfig(name)
	if err != nil {
		return nil, err
	}
	return r.Store.Bucket(info.Name)
}

// List returns the list of all buckets
func (r *Registry) List() ([]string, error) {
	var buckets []internal.Bucket

	err := r.DB.All(&buckets)
	if err != nil {
		return nil, errors.Wrap(err, "boltdb.registry.List failed to fetch buckets")
	}

	names := make([]string, len(buckets))
	for i := range buckets {
		names[i] = buckets[i].Name
	}

	return names, nil
}

// Close BoltDB connection
func (r *Registry) Close() error {
	return r.DB.Close()
}
