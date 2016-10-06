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
	return c.App.Registry.Create(name)
}

func (c *cli) Save(bucketName, key string, data []byte) error {
	bucket, err := store.GetBucketOrCreate(c.App.Registry, c.App.Store, bucketName)
	if err != nil {
		return err
	}
	defer bucket.Close()

	data = json.ToValidJSON(data)

	_, err = bucket.Save(key, data)
	return err
}

func (c *cli) Get(bucketName, key string) ([]byte, error) {
	info, err := c.App.Registry.BucketInfo(bucketName)
	if err != nil {
		return nil, err
	}
	bucket, err := c.App.Store.Bucket(info.Name)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	item, err := bucket.Get(key)
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (c *cli) List(bucketName string) ([]brazier.Item, error) {
	info, err := c.App.Registry.BucketInfo(bucketName)
	if err != nil {
		return nil, err
	}
	bucket, err := c.App.Store.Bucket(info.Name)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	return bucket.Page(1, -1)
}

func (c *cli) ListBuckets() ([]string, error) {
	return c.App.Registry.List()
}

func (c *cli) Delete(bucket, key string) error {
	b, err := c.App.Store.Bucket(bucket)
	if err != nil {
		return err
	}
	defer b.Close()

	return b.Delete(key)
}
