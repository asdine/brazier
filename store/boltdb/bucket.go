package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

// NewBucket returns a Bucket
func NewBucket(s *Store, key string, node storm.Node) *Bucket {
	return &Bucket{
		key:   key,
		store: s,
		node:  node,
	}
}

// Bucket is a BoltDB implementation a bucket
type Bucket struct {
	key   string
	store *Store
	node  storm.Node
}

// Save user data to the bucket. Returns an Iten
func (b *Bucket) Save(key string, data []byte) (*brazier.Item, error) {
	err := b.node.Set("items", key, data)
	if err != nil {
		return nil, err
	}

	return &brazier.Item{
		Key:  key,
		Data: data,
	}, nil
}

// Get an item by id
func (b *Bucket) Get(key string) (*brazier.Item, error) {
	var data []byte
	err := b.node.Get("items", key, &data)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.bucket.Get failed to fetch item")
	}

	return &brazier.Item{
		Key:  key,
		Data: data,
	}, nil
}

// Delete item from the bucket
func (b *Bucket) Delete(key string) error {
	return b.node.Delete("items", key)
}

// Page returns a list of items
func (b *Bucket) Page(page int, perPage int) ([]brazier.Item, error) {
	var skip int

	if page <= 0 {
		return nil, nil
	}

	if perPage >= 0 {
		skip = (page - 1) * perPage
	}

	var items []brazier.Item
	err := b.node.Select().Bucket("items").Skip(skip).Limit(perPage).RawEach(func(k, v []byte) error {
		items = append(items, brazier.Item{
			Key:  string(k),
			Data: v,
		})

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "boltdb.bucket.Page failed to fetch items")
	}

	return items, nil
}

// Close the bucket session
func (b *Bucket) Close() error {
	return b.store.closeSession(b.key)
}
