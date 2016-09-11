package rpc_test

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/stretchr/testify/require"
)

func TestGetter(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := internal.NewGetterClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false
	item, err := b.Save("key", []byte("data"))
	require.NoError(t, err)
	b.SaveInvoked = false

	r, err := c.Get(context.Background(), &internal.GetRequest{Bucket: "bucket", Key: "key"})
	require.NoError(t, err)
	require.Equal(t, "key", r.Key)
	require.True(t, s.BucketInvoked)
	require.True(t, b.GetInvoked)
	require.Equal(t, item.Data, r.Data)

	r, err = c.Get(context.Background(), &internal.GetRequest{Bucket: "bucket", Key: "unknown key"})
	require.Error(t, err)
}
