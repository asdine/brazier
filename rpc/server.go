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
func NewServer(s *store.Store) brazier.Server {
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

// Server is the Brazier gRPC server.
type Server struct {
	Store *store.Store
}

// Create a bucket.
func (s *Server) Create(ctx context.Context, in *proto.Selector) (*proto.Empty, error) {
	err := s.Store.CreateBucket(in.Path)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Put an item in the bucket.
func (s *Server) Put(ctx context.Context, in *proto.NewItem) (*proto.Empty, error) {
	data := json.ToValidJSON(in.Value)

	_, err := s.Store.Put(in.Path, data)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// Get an item from the bucket.
func (s *Server) Get(ctx context.Context, in *proto.Selector) (*proto.Item, error) {
	item, err := s.Store.Get(in.Path)
	if err != nil {
		return nil, err
	}

	r := proto.Item{
		Key:   item.Key,
		Value: item.Data,
	}

	return &r, nil
}

// Delete an item from the bucket.
func (s *Server) Delete(ctx context.Context, in *proto.Selector) (*proto.Empty, error) {
	err := s.Store.Delete(in.Path)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// List the content of a bucket.
func (s *Server) List(ctx context.Context, in *proto.Selector) (*proto.Tree, error) {
	var items []brazier.Item
	var err error

	if in.Recursive {
		items, err = s.Store.Tree(in.Path)
	} else {
		items, err = s.Store.List(in.Path, 1, -1)
	}
	if err != nil {
		return nil, err
	}

	return &proto.Tree{Children: s.tree(items)}, nil
}

func (s *Server) tree(items []brazier.Item) []*proto.Node {
	list := make([]*proto.Node, len(items))
	for i := range items {
		list[i] = &proto.Node{
			Key:   items[i].Key,
			Value: items[i].Data,
		}

		if items[i].Children != nil {
			list[i].Children = s.tree(items[i].Children)
		}
	}

	return list
}
