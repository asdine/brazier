package rpc

import (
	"net"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewServer returns a configured gRPC server
func NewServer(s brazier.Store) brazier.Server {
	g := grpc.NewServer()
	srv := Server{Store: s}
	proto.RegisterBucketServer(g, &srv)
	return &serverWrapper{srv: g}
}

type serverWrapper struct {
	srv *grpc.Server
}

func (s *serverWrapper) Serve(l net.Listener) error {
	return s.srv.Serve(l)
}

func (s *serverWrapper) Stop(time.Duration) {
	s.srv.GracefulStop()
}

// Server is the Brazier gRPC server
type Server struct {
	Store brazier.Store
}

// Create a bucket
func (s *Server) Create(ctx context.Context, in *proto.NewBucket) (*proto.Empty, error) {
	err := s.Store.Create(in.Name)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Buckets returns the list of existing buckets
func (s *Server) Buckets(ctx context.Context, in *proto.Empty) (*proto.BucketInfos, error) {
	list, err := s.Store.List()
	if err != nil {
		return nil, err
	}

	var infos proto.BucketInfos
	infos.Buckets = make([]*proto.BucketInfo, len(list))
	for i, name := range list {
		infos.Buckets[i] = &proto.BucketInfo{
			Name: name,
		}
	}

	return &infos, nil
}

// Save an item to the bucket
func (s *Server) Save(ctx context.Context, in *proto.NewItem) (*proto.Empty, error) {
	bucket, err := store.GetBucketOrCreate(s.Store, in.Bucket)
	if err != nil {
		return nil, err
	}

	data := json.ToValidJSON(in.Data)
	_, err = bucket.Save(in.Key, data)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Get an item from the bucket
func (s *Server) Get(ctx context.Context, in *proto.KeySelector) (*proto.Item, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	item, err := b.Get(in.Key)
	if err != nil {
		return nil, err
	}

	r := proto.Item{
		Key:  item.Key,
		Data: item.Data,
	}

	return &r, nil
}

// Delete an item from the bucket
func (s *Server) Delete(ctx context.Context, in *proto.KeySelector) (*proto.Empty, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	err = b.Delete(in.Key)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// List the content of a bucket
func (s *Server) List(ctx context.Context, in *proto.BucketSelector) (*proto.Items, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	items, err := b.Page(1, -1)
	if err != nil {
		return nil, err
	}

	list := make([]*proto.Item, len(items))
	for i := range items {
		list[i] = &proto.Item{
			Key:  items[i].Key,
			Data: items[i].Data,
		}
	}

	return &proto.Items{Items: list}, nil
}
