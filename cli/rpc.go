package cli

import (
	"golang.org/x/net/context"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc/proto"
)

type rpcCli struct {
	App    *app
	Client proto.BucketClient
}

func (r *rpcCli) Create(path string) error {
	_, err := r.Client.Create(context.Background(), &proto.Selector{Path: path})
	return err
}

func (r *rpcCli) Save(path string, data []byte) error {
	_, err := r.Client.Save(context.Background(), &proto.NewItem{Path: path, Data: data})
	return err
}

func (r *rpcCli) Get(path string) ([]byte, error) {
	item, err := r.Client.Get(context.Background(), &proto.Selector{Path: path})
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (r *rpcCli) List(path string) ([]brazier.Item, error) {
	resp, err := r.Client.List(context.Background(), &proto.Selector{Path: path})
	if err != nil {
		return nil, err
	}

	items := make([]brazier.Item, len(resp.Items))
	for i, item := range resp.Items {
		items[i].Key = item.Key
		items[i].Data = item.Data
	}

	return items, nil
}

func (r *rpcCli) Delete(path string) error {
	_, err := r.Client.Delete(context.Background(), &proto.Selector{Path: path})
	return err
}
