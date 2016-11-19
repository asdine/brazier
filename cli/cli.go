package cli

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
)

// Cli handles command line requests
type Cli interface {
	Create(path string) error
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	List(path string, recursive bool) ([]brazier.Item, error)
	Delete(path string) error
}

type cli struct {
	App *app
}

func (c *cli) Create(path string) error {
	return c.App.Store.CreateBucket(path)
}

func (c *cli) Put(path string, data []byte) error {
	data = json.ToValidJSON(data)

	_, err := c.App.Store.Put(path, data)
	return err
}

func (c *cli) Get(path string) ([]byte, error) {
	item, err := c.App.Store.Get(path)
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (c *cli) List(path string, recursive bool) ([]brazier.Item, error) {
	if recursive {
		return c.App.Store.Tree(path)
	}

	return c.App.Store.List(path, 1, -1)
}

func (c *cli) Delete(path string) error {
	return c.App.Store.Delete(path)
}
