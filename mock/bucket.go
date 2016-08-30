package mock

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewBucket returns a Bucket
func NewBucket() *Bucket {
	return &Bucket{
		data: make(map[string][]byte),
	}
}

// Bucket is a mock implementation of a bucket
type Bucket struct {
	data          map[string][]byte
	SaveInvoked   bool
	GetInvoked    bool
	DeleteInvoked bool
	CloseInvoked  bool
}

// Save user data to the bucket. Returns an Iten
func (b *Bucket) Save(id string, data []byte) (*brazier.Item, error) {
	b.SaveInvoked = true

	b.data[id] = data
	return &brazier.Item{
		ID:        id,
		Data:      data,
		CreatedAt: time.Now(),
	}, nil
}

// Get an item by id
func (b *Bucket) Get(id string) (*brazier.Item, error) {
	b.GetInvoked = true

	if data, ok := b.data[id]; ok {
		return &brazier.Item{
			ID:   id,
			Data: data,
		}, nil
	}

	return nil, store.ErrNotFound
}

// Delete item from the bucket
func (b *Bucket) Delete(id string) error {
	b.DeleteInvoked = true

	if _, ok := b.data[id]; ok {
		delete(b.data, id)
		return nil
	}

	return store.ErrNotFound
}

// Close bucket
func (b *Bucket) Close() error {
	b.CloseInvoked = true
	return nil
}
