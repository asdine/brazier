package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

// NewBucket returns a Bucket
func NewBucket(s *Store, id string, node storm.Node) *Bucket {
	return &Bucket{
		id:    id,
		store: s,
		node:  node,
	}
}

// Bucket is a BoltDB implementation a bucket
type Bucket struct {
	id    string
	store *Store
	node  storm.Node
}

// Save user data to the bucket. Returns an Iten
func (b *Bucket) Save(id string, data []byte) (*brazier.Item, error) {
	var i item

	tx, err := b.node.Begin(true)
	if err != nil {
		return nil, err
	}

	err = tx.One("ID", id, &i)
	if err != nil {
		if err != storm.ErrNotFound {
			tx.Rollback()
			return nil, err
		}

		i = item{
			ID:        id,
			Data:      data,
			CreatedAt: time.Now(),
		}
	} else {
		i.UpdatedAt = time.Now()
		i.Data = data
	}

	err = tx.Save(&i)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &brazier.Item{
		ID:        i.ID,
		Data:      i.Data,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}, tx.Commit()
}

// Get an item by id
func (b *Bucket) Get(id string) (*brazier.Item, error) {
	var i item

	err := b.node.One("ID", id, &i)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.bucket.Get failed to fetch item")
	}

	return &brazier.Item{
		ID:        i.ID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		Data:      i.Data,
	}, nil
}

// Delete item from the bucket
func (b *Bucket) Delete(id string) error {
	var i item

	tx, err := b.node.Begin(true)
	if err != nil {
		return errors.Wrap(err, "boltdb.bucket.Delete failed to create transaction")
	}

	err = tx.One("ID", id, &i)
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

// Close the bucket session
func (b *Bucket) Close() error {
	return b.store.closeSession(b.id)
}
