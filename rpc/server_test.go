package rpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newServer(t *testing.T, s brazier.Store) (*grpc.ClientConn, func()) {
	l, err := net.Listen("tcp", ":")
	require.NoError(t, err)

	srv := grpc.NewServer()
	bSrv := rpc.Server{Store: s}
	proto.RegisterBucketServer(srv, &bSrv)

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

func TestCreator(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Create(context.Background(), &proto.NewBucket{Name: "bucket"})
	require.NoError(t, err)

	require.True(t, s.CreateInvoked)
	_, err = s.Bucket("bucket")
	require.NoError(t, err)
}

func TestSaver(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Save(context.Background(), &proto.NewItem{Bucket: "bucket", Key: "key", Data: []byte("data")})
	require.NoError(t, err)

	require.True(t, s.BucketInvoked)
	bucket, err := s.Bucket("bucket")
	b := bucket.(*mock.Bucket)
	require.NoError(t, err)
	require.True(t, b.SaveInvoked)

	item, err := b.Get("key")
	require.NoError(t, err)
	require.Equal(t, []byte("data"), item.Data)
}

func TestLister(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := proto.NewBucketClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false

	r, err := c.List(context.Background(), &proto.BucketSelector{Bucket: "bucket"})
	require.NoError(t, err)
	require.Len(t, r.Items, 0)
	require.True(t, s.BucketInvoked)
	require.True(t, b.PageInvoked)

	s.BucketInvoked = false
	b.PageInvoked = false

	var list [][]byte

	for i := 0; i < 20; i++ {
		item, err := b.Save(fmt.Sprintf("key%d", i), []byte("data"))
		require.NoError(t, err)
		b.SaveInvoked = false
		list = append(list, item.Data)
	}

	r, err = c.List(context.Background(), &proto.BucketSelector{Bucket: "bucket"})
	require.NoError(t, err)
	for i := 0; i < 20; i++ {
		require.Equal(t, list[i], r.Items[i].Data)
	}
	require.True(t, s.BucketInvoked)
	require.True(t, b.PageInvoked)
}

func TestGetter(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := proto.NewBucketClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false
	item, err := b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	r, err := c.Get(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.Equal(t, "key", r.Key)
	require.True(t, s.BucketInvoked)
	require.True(t, b.GetInvoked)
	require.Equal(t, item.Data, r.Data)

	r, err = c.Get(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "unknown key"})
	require.Error(t, err)
}

func TestDeleter(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := proto.NewBucketClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false
	_, err = b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	_, err = c.Delete(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.True(t, s.BucketInvoked)
	require.True(t, b.DeleteInvoked)

	_, err = b.Get("key")
	require.Equal(t, store.ErrNotFound, err)

	_, err = c.Delete(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.Error(t, err)
}
