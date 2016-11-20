package rpc_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func newServer(t *testing.T, s *store.Store) (*grpc.ClientConn, func()) {
	l, err := net.Listen("tcp", ":")
	require.NoError(t, err)

	srv := rpc.NewServer(s)

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

func TestCreate(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(r)
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Create(context.Background(), &proto.Selector{Path: "a/b/c/"})
	require.NoError(t, err)

	require.True(t, r.CreateInvoked)
	_, err = r.Bucket("a", "b", "c")
	require.NoError(t, err)
}

func TestPut(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(r)
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	_, err := c.Put(context.Background(), &proto.NewItem{Path: "a/b/c", Value: []byte("data")})
	require.NoError(t, err)

	require.True(t, r.BucketInvoked)
	bucket, err := r.Bucket("a", "b")
	b := bucket.(*mock.Bucket)
	require.NoError(t, err)
	require.True(t, b.SaveInvoked)

	item, err := b.Get("c")
	require.NoError(t, err)
	require.Equal(t, []byte(`"data"`), item.Data)
}

func TestList(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(r)
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("a", "b", "c")
	require.NoError(t, err)
	bucket, err := r.Bucket("a", "b", "c")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false

	t.Run("list", func(t *testing.T) {
		resp, err := c.List(context.Background(), &proto.Selector{Path: "a/b/c/"})
		require.NoError(t, err)
		require.Len(t, resp.Children, 0)
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

		resp, err = c.List(context.Background(), &proto.Selector{Path: "a/b/c/"})
		require.NoError(t, err)
		for i := 0; i < 20; i++ {
			require.Equal(t, list[i], resp.Children[i].Value)
		}
		require.True(t, r.BucketInvoked)
		require.True(t, b.PageInvoked)

		_, err = c.List(context.Background(), &proto.Selector{Path: "a/b/d"})
		require.Error(t, err)
	})

	t.Run("tree", func(t *testing.T) {
		resp, err := c.List(context.Background(), &proto.Selector{Path: "a/", Recursive: true})
		require.NoError(t, err)
		require.Len(t, resp.Children, 1)
		require.Equal(t, "b", resp.Children[0].Key)
		require.Len(t, resp.Children[0].Children, 1)
		require.Equal(t, "c", resp.Children[0].Children[0].Key)
		require.Len(t, resp.Children[0].Children[0].Children, 20)
		require.Equal(t, "key0", resp.Children[0].Children[0].Children[0].Key)
	})
}

func TestGet(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(r)
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("a", "b")
	require.NoError(t, err)
	bucket, err := r.Bucket("a", "b")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false
	item, err := b.Save("c", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	resp, err := c.Get(context.Background(), &proto.Selector{Path: "a/b/c"})
	require.NoError(t, err)
	require.Equal(t, "c", resp.Key)
	require.True(t, r.BucketInvoked)
	require.True(t, b.GetInvoked)
	require.Equal(t, item.Data, resp.Value)

	resp, err = c.Get(context.Background(), &proto.Selector{Path: "a/b/d"})
	require.Error(t, err)
}

func TestDelete(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(r)
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := proto.NewBucketClient(conn)

	err := r.Create("a", "b")
	require.NoError(t, err)
	bucket, err := r.Bucket("a", "b")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	r.BucketInvoked = false
	_, err = b.Save("c", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	_, err = c.Delete(context.Background(), &proto.Selector{Path: "a/b/c"})
	require.NoError(t, err)
	require.True(t, r.BucketInvoked)
	require.True(t, b.DeleteInvoked)

	_, err = b.Get("c")
	require.Equal(t, store.ErrNotFound, err)

	_, err = c.Delete(context.Background(), &proto.Selector{Path: "a/b/d"})
	require.Error(t, err)
}
