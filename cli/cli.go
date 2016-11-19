package cli

import (
	"strings"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
)

// Cli handles command line requests
type Cli interface {
	Create(path string) error
	Put(path string, data []byte) error
	Get(path string, recursive bool) ([]byte, error)
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

func (c *cli) Get(path string, recursive bool) ([]byte, error) {
	if strings.HasSuffix(path, "/") {
		var items []brazier.Item
		var err error

		if recursive {
			items, err = c.App.Store.Tree(path)
		} else {
			items, err = c.App.Store.List(path, 1, -1)
		}

		if err != nil {
			return nil, err
		}

		data, err := json.MarshalList(items)
		if err != nil {
			return nil, err
		}

		return append(data, '\n'), nil
	}

	item, err := c.App.Store.Get(path)
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (c *cli) Delete(path string) error {
	return c.App.Store.Delete(path)
}
