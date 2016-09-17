package cli

import (
	"bytes"

	"golang.org/x/net/context"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/asdine/brazier/rpc/proto"
)

type rpcCli struct {
	App    *app
	Client proto.BucketClient
}

func (r *rpcCli) Create(name string) error {
	_, err := r.Client.Create(context.Background(), &proto.NewBucket{Name: name})
	return err
}

func (r *rpcCli) Save(bucket, key string, data []byte) error {
	_, err := r.Client.Save(context.Background(), &proto.NewItem{Bucket: bucket, Key: key, Data: data})
	return err
}

func (r *rpcCli) Get(bucket, key string) ([]byte, error) {
	item, err := r.Client.Get(context.Background(), &proto.KeySelector{Bucket: bucket, Key: key})
	if err != nil {
		return nil, err
	}

	return append(item.Data, '\n'), nil
}

func (r *rpcCli) List(bucket string) ([]byte, error) {
	resp, err := r.Client.List(context.Background(), &proto.BucketSelector{Bucket: bucket})
	if err != nil {
		return nil, err
	}

	items := make([]brazier.Item, len(resp.Items))
	for i, item := range resp.Items {
		items[i].Key = item.Key
		items[i].Data = item.Data
	}

	data, err := json.MarshalList(items)
	if err != nil {
		return nil, err
	}

	return append(data, '\n'), nil
}

func (r *rpcCli) ListBuckets() ([]byte, error) {
	resp, err := r.Client.Buckets(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	for i := range resp.Buckets {
		buffer.WriteString(resp.Buckets[i].Name)
		buffer.WriteByte('\n')
	}

	return buffer.Bytes(), nil
}

func (r *rpcCli) Delete(bucket, key string) error {
	_, err := r.Client.Delete(context.Background(), &proto.KeySelector{Bucket: bucket, Key: key})
	return err
}
