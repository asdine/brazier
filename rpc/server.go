package rpc

import (
	"net"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/asdine/brazier/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewServer returns a configured gRPC server
func NewServer(s brazier.Store) brazier.Server {
	g := grpc.NewServer()
	srv := Server{Store: s}
	internal.RegisterSaverServer(g, &srv)
	internal.RegisterGetterServer(g, &srv)
	internal.RegisterListerServer(g, &srv)
	internal.RegisterDeleterServer(g, &srv)
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

// Save an item to the bucket
func (s *Server) Save(ctx context.Context, in *internal.SaveRequest) (*internal.SaveReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}
		err = s.Store.Create(in.Bucket)
		if err != nil {
			return nil, err
		}
		b, err = s.Store.Bucket(in.Bucket)
		if err != nil {
			return nil, err
		}
	}

	_, err = b.Save(in.Key, in.Data)
	if err != nil {
		return nil, err
	}

	return &internal.SaveReply{Status: 200}, nil
}

// Get an item from the bucket
func (s *Server) Get(ctx context.Context, in *internal.GetRequest) (*internal.GetReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	item, err := b.Get(in.Key)
	if err != nil {
		return nil, err
	}

	r := internal.GetReply{
		Key:  item.Key,
		Data: item.Data,
	}

	return &r, nil
}

// Delete an item from the bucket
func (s *Server) Delete(ctx context.Context, in *internal.DeleteRequest) (*internal.DeleteReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	err = b.Delete(in.Key)
	if err != nil {
		return nil, err
	}

	return &internal.DeleteReply{Status: 200}, nil
}

// List the content of a bucket
func (s *Server) List(ctx context.Context, in *internal.ListRequest) (*internal.ListReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	items, err := b.Page(1, -1)
	if err != nil {
		return nil, err
	}

	list := make([][]byte, len(items))
	for i := range items {
		list[i] = items[i].Data
	}

	return &internal.ListReply{Items: list}, nil
}
