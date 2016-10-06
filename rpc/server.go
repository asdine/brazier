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
func NewServer(r brazier.Registry, s brazier.Store) brazier.Server {
	g := grpc.NewServer()
	srv := Server{Registry: r, Store: s}
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
	Registry brazier.Registry
	Store    brazier.Store
}

// Create a bucket
func (s *Server) Create(ctx context.Context, in *proto.NewBucket) (*proto.Empty, error) {
	err := s.Registry.Create(in.Name)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Buckets returns the list of existing buckets
func (s *Server) Buckets(ctx context.Context, in *proto.Empty) (*proto.BucketInfos, error) {
	list, err := s.Registry.List()
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
	bucket, err := store.GetBucketOrCreate(s.Registry, s.Store, in.Bucket)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	data := json.ToValidJSON(in.Data)
	_, err = bucket.Save(in.Key, data)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Get an item from the bucket
func (s *Server) Get(ctx context.Context, in *proto.KeySelector) (*proto.Item, error) {
	info, err := s.Registry.BucketInfo(in.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := s.Store.Bucket(info.Name)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	item, err := bucket.Get(in.Key)
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
	info, err := s.Registry.BucketInfo(in.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := s.Store.Bucket(info.Name)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	err = bucket.Delete(in.Key)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// List the content of a bucket
func (s *Server) List(ctx context.Context, in *proto.BucketSelector) (*proto.Items, error) {
	info, err := s.Registry.BucketInfo(in.Bucket)
	if err != nil {
		return nil, err
	}
	bucket, err := s.Store.Bucket(info.Name)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	items, err := bucket.Page(1, -1)
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
