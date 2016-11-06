package boltdb

import (
	"bytes"
	"strings"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// NewRegistry returns a BoltDB Registry.
func NewRegistry(path string, b brazier.Backend) (*Registry, error) {
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
		DB:      db,
		Backend: b,
	}, nil
}

// Registry is a BoltDB registry.
type Registry struct {
	DB      *storm.DB
	Backend brazier.Backend
}

// Create a bucket in the registry.
func (r *Registry) Create(nodes ...string) error {
	path := strings.Join(nodes, "/")

	err := r.DB.Bolt.Update(func(tx *bolt.Tx) error {
		b := r.DB.GetBucket(tx, nodes...)
		if b != nil {
			return storm.ErrAlreadyExists
		}

		var n storm.Node
		n = r.DB

		last := nodes[len(nodes)-1]
		if len(nodes) > 1 {
			nodes = nodes[:len(nodes)-1]
			n = r.DB.From(nodes...)
		}

		n = n.WithTransaction(tx)
		_, err := n.CreateBucketIfNotExists(tx, last)
		if err != nil {
			return errors.Wrapf(err, "failed to create bucket at path %s", path)
		}

		return nil
	})

	if err == storm.ErrAlreadyExists {
		return store.ErrAlreadyExists
	}

	return errors.Wrapf(err, "failed to create bucket at path %s", path)
}

// Bucket returns the selected bucket from the Backend.
func (r *Registry) Bucket(nodes ...string) (brazier.Bucket, error) {
	err := r.DB.Bolt.View(func(tx *bolt.Tx) error {
		b := r.DB.GetBucket(tx, nodes...)
		if b == nil {
			return storm.ErrNotFound
		}

		return nil
	})

	if err == storm.ErrNotFound {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch bucket at path %s", strings.Join(nodes, "/"))
	}

	return r.Backend.Bucket(nodes...)
}

// Children buckets of the specified path.
func (r *Registry) Children(nodes ...string) ([]brazier.Item, error) {
	var items []brazier.Item
	var err error

	err = r.DB.Bolt.View(func(tx *bolt.Tx) error {
		b := r.DB.GetBucket(tx, nodes...)
		if b == nil {
			return store.ErrNotFound
		}

		items, err = r.childrenOf(b)

		return err
	})

	return items, err
}

func (r *Registry) childrenOf(b *bolt.Bucket) ([]brazier.Item, error) {
	var items []brazier.Item
	var err error

	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil || bytes.HasPrefix(k, []byte("__storm")) {
			continue
		}

		item := brazier.Item{
			Key: string(k),
		}

		item.Children, err = r.childrenOf(b.Bucket(k))
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// Close BoltDB connection
func (r *Registry) Close() error {
	err := r.Backend.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close backend")
	}

	err = r.DB.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close registry")
	}

	return nil
}
