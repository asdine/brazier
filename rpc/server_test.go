package rpc_test

import (
	"net"
	"testing"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newServer(t *testing.T, s brazier.Store) (*grpc.ClientConn, func()) {
	l, err := net.Listen("tcp", ":")
	require.NoError(t, err)

	srv := grpc.NewServer()
	bSrv := rpc.Server{Store: s}
	internal.RegisterSaverServer(srv, &bSrv)
	internal.RegisterGetterServer(srv, &bSrv)
	internal.RegisterDeleterServer(srv, &bSrv)
	internal.RegisterListerServer(srv, &bSrv)

	go func() {
		srv.Serve(l)
	}()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
		time.Sleep(5 * time.Millisecond)
		srv.Stop()
	}
}
