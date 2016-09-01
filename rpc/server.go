package rpc

import (
	"fmt"
	"net"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server is the Brazier gRPC server
type Server struct {
	Store brazier.Store
}

// Save an item to the store
func (s *Server) Save(ctx context.Context, in *proto.SaveRequest) (*proto.SaveReply, error) {
	b, err := s.Store.Bucket(in.Bucket)
	if err != nil {
		return nil, err
	}

	_, err = b.Save(in.Key, in.Data)
	if err != nil {
		return nil, err
	}

	return &proto.SaveReply{Status: 200}, nil
}

// Get an item from the store
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
