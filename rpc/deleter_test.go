package rpc_test

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func TestDeleter(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := internal.NewDeleterClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false
	_, err = b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	r, err := c.Delete(context.Background(), &internal.DeleteRequest{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.Equal(t, int32(200), r.Status)
	require.True(t, s.BucketInvoked)
	require.True(t, b.DeleteInvoked)

	_, err = b.Get("key")
	require.Equal(t, store.ErrNotFound, err)

	r, err = c.Delete(context.Background(), &internal.DeleteRequest{Bucket: "bucket", Key: "key"})
	require.Error(t, err)
}
