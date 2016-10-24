package rpc_test

import (
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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func newServer(t *testing.T, r brazier.Registry) (*grpc.ClientConn, func()) {
	l, err := net.Listen("tcp", ":")
	require.NoError(t, err)

	srv := rpc.NewServer(r)

	go func() {
		srv.Serve(l)
	}()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
		time.Sleep(5 * time.Millisecond)
		srv.Stop(1 * time.Second)
	}
}

func TestCreator(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Create(context.Background(), &proto.NewBucket{Name: "bucket"})
	require.NoError(t, err)

	require.True(t, r.CreateInvoked)
	_, err = r.Bucket("bucket")
	require.NoError(t, err)
}

func TestBuckets(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	list, err := c.Buckets(context.Background(), &proto.Empty{})
	require.NoError(t, err)
	require.True(t, r.ListInvoked)
	require.Len(t, list.Buckets, 0)

	r.Create("bucket1")
	r.Create("bucket2")

	list, err = c.Buckets(context.Background(), &proto.Empty{})
	require.NoError(t, err)
	require.True(t, r.ListInvoked)
	require.Len(t, list.Buckets, 2)
	require.NoError(t, err)
}

func TestSaver(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Save(context.Background(), &proto.NewItem{Bucket: "bucket", Key: "key", Data: []byte("data")})
	require.NoError(t, err)

	require.True(t, r.BucketInvoked)
	bucket, err := r.Bucket("bucket")
	b := bucket.(*mock.Bucket)
	require.NoError(t, err)
	require.True(t, b.SaveInvoked)

	item, err := b.Get("key")
	require.NoError(t, err)
	require.Equal(t, []byte(`"data"`), item.Data)
}

func TestLister(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("bucket")
	require.NoError(t, err)
	bucket, err := r.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false

	resp, err := c.List(context.Background(), &proto.BucketSelector{Bucket: "bucket"})
	require.NoError(t, err)
	require.Len(t, resp.Items, 0)
	require.True(t, r.BucketInvoked)
	require.True(t, b.PageInvoked)

	r.BucketInvoked = false
	b.PageInvoked = false

	var list [][]byte

	for i := 0; i < 20; i++ {
		item, err := b.Save(fmt.Sprintf("key%d", i), []byte("data"))
		require.NoError(t, err)
		b.SaveInvoked = false
		list = append(list, item.Data)
	}

	resp, err = c.List(context.Background(), &proto.BucketSelector{Bucket: "bucket"})
	require.NoError(t, err)
	for i := 0; i < 20; i++ {
		require.Equal(t, list[i], resp.Items[i].Data)
	}
	require.True(t, r.BucketInvoked)
	require.True(t, b.PageInvoked)

	resp, err = c.List(context.Background(), &proto.BucketSelector{Bucket: "something"})
	require.NoError(t, err)
	require.Len(t, resp.Items, 0)
}

func TestGetter(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("bucket")
	require.NoError(t, err)
	bucket, err := r.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false
	item, err := b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	resp, err := c.Get(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.Equal(t, "key", resp.Key)
	require.True(t, r.BucketInvoked)
	require.True(t, b.GetInvoked)
	require.Equal(t, item.Data, resp.Data)

	resp, err = c.Get(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "unknown key"})
	require.Error(t, err)
}

func TestDeleter(t *testing.T) {
	r := mock.NewRegistry(mock.NewStore())
	conn, cleanup := newServer(t, r)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("bucket")
	require.NoError(t, err)
	bucket, err := r.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false
	_, err = b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	_, err = c.Delete(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.True(t, r.BucketInvoked)
	require.True(t, b.DeleteInvoked)

	_, err = b.Get("key")
	require.Equal(t, store.ErrNotFound, err)

	_, err = c.Delete(context.Background(), &proto.KeySelector{Bucket: "bucket", Key: "key"})
	require.Error(t, err)
}
