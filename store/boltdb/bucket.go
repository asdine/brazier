package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb/internal"
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
	var i internal.Item

	tx, err := b.node.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.One("ID", key, &i)
	if err != nil {
		if err != storm.ErrNotFound {
			return nil, err
		}

		i = internal.Item{
			ID:        key,
			Data:      data,
			CreatedAt: time.Now().UnixNano(),
		}
	} else {
		i.UpdatedAt = time.Now().UnixNano()
		i.Data = data
	}

	err = tx.Save(&i)
	if err != nil {
		return nil, err
	}

	return &brazier.Item{
		Key:  i.ID,
		Data: i.Data,
	}, tx.Commit()
}

// Get an item by key
func (b *Bucket) Get(key string) (*brazier.Item, error) {
	var i internal.Item

	err := b.node.One("ID", key, &i)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.bucket.Get failed to fetch item")
	}

	return &brazier.Item{
		Key:  i.ID,
		Data: i.Data,
	}, nil
}

// Delete item from the bucket
func (b *Bucket) Delete(key string) error {
	var i internal.Item

	tx, err := b.node.Begin(true)
	if err != nil {
		return errors.Wrap(err, "boltdb.bucket.Delete failed to create transaction")
	}

	err = tx.One("ID", key, &i)
	if err != nil {
		tx.Rollback()
		if err == storm.ErrNotFound {
			return store.ErrNotFound
		}
		return errors.Wrap(err, "boltdb.bucket.Delete failed to fetch item")
	}

	err = tx.DeleteStruct(&i)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "boltdb.bucket.Delete failed to delete item")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "boltdb.bucket.Delete failed to commit")
	}

	return nil
}

// Page returns a list of items
func (b *Bucket) Page(page int, perPage int) ([]brazier.Item, error) {
	var skip int
	var list []internal.Item

	if page <= 0 {
		return nil, nil
	}

	if perPage >= 0 {
		skip = (page - 1) * perPage
	}

	err := b.node.All(&list, storm.Skip(skip), storm.Limit(perPage))
	if err != nil {
		return nil, errors.Wrap(err, "boltdb.bucket.Page failed to fetch items")
	}

	items := make([]brazier.Item, len(list))
	for i := range list {
		items[i].Key = list[i].ID
		items[i].Data = list[i].Data
	}
	return items, nil
}

// Close the bucket session
func (b *Bucket) Close() error {
	return b.store.closeSession(b.key)
}
