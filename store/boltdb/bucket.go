package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/storm"
	"github.com/dchest/uniuri"
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

// Close the session of the bucket
func (b *Bucket) Close() error {
	return nil
}
