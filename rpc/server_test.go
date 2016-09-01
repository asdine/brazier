package rpc_test

import (
	"net"
	"testing"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newServer(t *testing.T, s brazier.Store) (*grpc.ClientConn, func()) {
	l, err := net.Listen("tcp", ":")
	require.NoError(t, err)

	srv := grpc.NewServer()
	proto.RegisterSaverServer(srv, &rpc.Server{Store: s})
	go func() {
		srv.Serve(l)
	}()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
		srv.Stop()
	}
}
