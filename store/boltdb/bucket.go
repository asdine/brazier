package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

// NewBucket returns a Bucket
func NewBucket(node *storm.Node) *Bucket {
	return &Bucket{
		node: node,
	}
}

// Bucket is a BoltDB implementation a bucket
type Bucket struct {
	node *storm.Node
}

// Add user data to the bucket. Returns an Iten
func (b *Bucket) Add(data []byte, mimeType string, name string) (*brazier.Item, error) {
	i := item{
		Data:      data,
		MimeType:  mimeType,
		PublicID:  name,
		CreatedAt: time.Now(),
	}

	if i.PublicID == "" {
		i.PublicID = uniuri.NewLen(10)
	}

	err := b.node.Save(&i)
	if err != nil {
		return nil, err
	}

	return &brazier.Item{
		ID:        i.PublicID,
		Data:      i.Data,
		MimeType:  i.MimeType,
		CreatedAt: i.CreatedAt,
	}, nil
}

// Get an item by id
func (b *Bucket) Get(id string) (*brazier.Item, error) {
	var item item

	err := b.node.One("PublicID", id, &item)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.bucket.Get failed to fetch item")
	}

	return &brazier.Item{
		ID:        item.PublicID,
		CreatedAt: item.CreatedAt,
		Data:      item.Data,
		MimeType:  item.MimeType,
	}, nil
}

// Delete item from the bucket
func (b *Bucket) Delete(id string) error {
	var i item

	tx, err := b.node.Begin(true)
	if err != nil {
		return errors.Wrap(err, "boltdb.bucket.Delete failed to create transaction")
	}

	err = tx.One("PublicID", id, &i)
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

// Close the session of the bucket
func (b *Bucket) Close() error {
	return nil
}
