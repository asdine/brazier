package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/stretchr/testify/require"
)

func TestLister(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()
	c := internal.NewListerClient(conn)

	err := s.Create("bucket")
	require.NoError(t, err)
	bucket, err := s.Bucket("bucket")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)
	s.BucketInvoked = false

	r, err := c.List(context.Background(), &internal.ListRequest{Bucket: "bucket"})
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

	r, err = c.List(context.Background(), &internal.ListRequest{Bucket: "bucket"})
	require.NoError(t, err)
	require.Equal(t, list, r.Items)
	require.True(t, s.BucketInvoked)
	require.True(t, b.PageInvoked)
}
