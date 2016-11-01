package mock

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// NewBucket returns a Bucket.
func NewBucket(name string) *Bucket {
	return &Bucket{
		Name: name,
		data: make(map[string]*brazier.Item),
	}
}

// Bucket is a mock implementation of a bucket.
type Bucket struct {
	Name          string
	data          map[string]*brazier.Item
	index         []*brazier.Item
	children      []*Bucket
	SaveInvoked   bool
	GetInvoked    bool
	DeleteInvoked bool
	PageInvoked   bool
	CloseInvoked  bool
}

// Save user data to the bucket. Returns an Item.
func (b *Bucket) Save(key string, data []byte) (*brazier.Item, error) {
	b.SaveInvoked = true

	item, ok := b.data[key]
	if !ok {
		item = &brazier.Item{
			Key:  key,
			Data: data,
		}
		b.data[key] = item
		b.index = append(b.index, item)
	} else {
		item.Data = data
	}

	return item, nil
}

// Get an item by key.
func (b *Bucket) Get(key string) (*brazier.Item, error) {
	b.GetInvoked = true

	if item, ok := b.data[key]; ok {
		return item, nil
	}

	return nil, store.ErrNotFound
}

// Delete item from the bucket
func (b *Bucket) Delete(key string) error {
	b.DeleteInvoked = true

	if _, ok := b.data[key]; ok {
		delete(b.data, key)
		return nil
	}

	return store.ErrNotFound
}

// Page returns a list of items.
func (b *Bucket) Page(page int, perPage int) ([]brazier.Item, error) {
	b.PageInvoked = true

	var start, end int

	if page <= 0 {
		return nil, nil
	}

	if perPage >= 0 {
		start = (page - 1) * perPage
	}

	if start >= len(b.index) {
		return nil, nil
	}

	if perPage == -1 {
		end = len(b.index)
	} else {
		end = start + perPage
		if end > len(b.index) {
			end = len(b.index)
		}
	}

	items := make([]brazier.Item, end-start)
	slice := b.index[start:end]
	for i := range slice {
		items[i] = *slice[i]
	}
	return items, nil
}

// Close bucket.
func (b *Bucket) Close() error {
	b.CloseInvoked = true
	return nil
}
