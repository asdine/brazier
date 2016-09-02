package mock

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewBucket returns a Bucket
func NewBucket() *Bucket {
	return &Bucket{
		data: make(map[string]*brazier.Item),
	}
}

// Bucket is a mock implementation of a bucket
type Bucket struct {
	data          map[string]*brazier.Item
	index         []*brazier.Item
	SaveInvoked   bool
	GetInvoked    bool
	DeleteInvoked bool
	PageInvoked   bool
	CloseInvoked  bool
}

// Save user data to the bucket. Returns an Iten
func (b *Bucket) Save(id string, data []byte) (*brazier.Item, error) {
	b.SaveInvoked = true

	item, ok := b.data[id]
	if !ok {
		item = &brazier.Item{
			ID:        id,
			Data:      data,
			CreatedAt: time.Now(),
		}
		b.data[id] = item
		b.index = append(b.index, item)
	} else {
		item.Data = data
		item.UpdatedAt = time.Now()
	}

	return item, nil
}

// Get an item by id
func (b *Bucket) Get(id string) (*brazier.Item, error) {
	b.GetInvoked = true

	if item, ok := b.data[id]; ok {
		return item, nil
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

// Page returns a list of items
func (b *Bucket) Page(page int, perPage int) ([]brazier.Item, error) {
	b.PageInvoked = true

	start := perPage * (page - 1)
	end := start + perPage
	if end > len(b.index) {
		end = len(b.index)
	}

	if start >= len(b.index) {
		return nil, nil
	}

	items := make([]brazier.Item, end-start)
	for i := range b.index[start:end] {
		items[i] = *b.index[i]
	}
	return items, nil
}

// Close bucket
func (b *Bucket) Close() error {
	b.CloseInvoked = true
	return nil
}
