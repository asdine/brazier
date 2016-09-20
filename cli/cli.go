package cli

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/asdine/brazier/store"
)

// Cli handles command line requests
type Cli interface {
	Create(name string) error
	Save(bucket, key string, data []byte) error
	Get(bucket, key string) ([]byte, error)
	List(bucket string) ([]brazier.Item, error)
	ListBuckets() ([]string, error)
	Delete(bucket, key string) error
}

type cli struct {
	App *app
}

func (c *cli) Create(name string) error {
	return c.App.Store.Create(name)
}

func (c *cli) Save(bucket, key string, data []byte) error {
	b, err := store.GetBucketOrCreate(c.App.Store, bucket)
	if err != nil {
		return err
	}
	defer b.Close()

	data = json.ToValidJSON(data)

	_, err = b.Save(key, data)
	return err
}

func (c *cli) Get(bucket, key string) ([]byte, error) {
	b, err := c.App.Store.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	defer b.Close()

	item, err := b.Get(key)
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (c *cli) List(bucket string) ([]brazier.Item, error) {
	b, err := c.App.Store.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	defer b.Close()

	return b.Page(1, -1)
}

func (c *cli) ListBuckets() ([]string, error) {
	return c.App.Store.List()
}

func (c *cli) Delete(bucket, key string) error {
	b, err := c.App.Store.Bucket(bucket)
	if err != nil {
		return err
	}
	defer b.Close()

	return b.Delete(key)
}
