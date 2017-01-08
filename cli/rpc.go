package cli

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
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

func (r *rpcCli) Put(path string, data []byte) error {
	_, err := r.Client.Put(context.Background(), &proto.NewItem{Path: path, Value: data})
	return err
}

func (r *rpcCli) Get(path string, recursive bool) ([]byte, error) {
	var err error
	var data []byte

	if strings.HasSuffix(path, "/") {
		resp, err := r.Client.List(context.Background(), &proto.Selector{Path: path, Recursive: recursive})
		if err != nil {
			return nil, err
		}

		data, err = json.MarshalListPretty(r.tree(resp.Children))
	} else {
		item, err := r.Client.Get(context.Background(), &proto.Selector{Path: path})
		if err != nil {
			return nil, err
		}

		data, err = json.PrettyPrintRaw(item.Value)
	}

	if err != nil {
		return nil, err
	}

	return append(data, '\n'), nil
}

func (r *rpcCli) tree(items []*proto.Node) []brazier.Item {
	list := make([]brazier.Item, len(items))
	for i, item := range items {
		list[i].Key = item.Key
		list[i].Data = item.Value

		if item.Children != nil {
			list[i].Children = r.tree(item.Children)
		}
	}

	return list
}

func (r *rpcCli) Delete(path string) error {
	_, err := r.Client.Delete(context.Background(), &proto.Selector{Path: path})
	return err
}
