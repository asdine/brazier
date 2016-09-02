package rpc

import (
	"fmt"
	"net"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server is the Brazier gRPC server
type Server struct {
	Store brazier.Store
}

// Save an item to the bucket
func (s *Server) Save(ctx context.Context, in *proto.SaveRequest) (*proto.SaveReply, error) {
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

	return &proto.SaveReply{Status: 200}, nil
}

// Get an item from the bucket
func (s *Server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	item, err := b.Get(in.Key)
	if err != nil {
		return nil, err
	}

	r := proto.GetReply{
		Key:       item.ID,
		CreatedAt: item.CreatedAt.UnixNano(),
		Data:      item.Data,
	}

	if !item.UpdatedAt.IsZero() {
		r.UpdatedAt = item.UpdatedAt.UnixNano()
	}

	return &r, nil
}

// Delete an item from the bucket
func (s *Server) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	err = b.Delete(in.Key)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteReply{Status: 200}, nil
}

// List the content of a bucket
func (s *Server) List(ctx context.Context, in *proto.ListRequest) (*proto.ListReply, error) {
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

	return &proto.ListReply{Items: list}, nil
}

// Serve runs the RPC server
func Serve(s brazier.Store, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	proto.RegisterSaverServer(srv, &Server{Store: s})
	proto.RegisterGetterServer(srv, &Server{Store: s})
	return srv.Serve(l)
}
