package rpc_test

import (
	"context"
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc/internal"
	"github.com/stretchr/testify/require"
)

func TestSaver(t *testing.T) {
	s := mock.NewStore()
	conn, cleanup := newServer(t, s)
	defer cleanup()

	c := internal.NewSaverClient(conn)

	r, err := c.Save(context.Background(), &internal.SaveRequest{Bucket: "bucket", Key: "key", Data: []byte("data")})
	require.NoError(t, err)
	require.Equal(t, int32(200), r.Status)

	require.True(t, s.BucketInvoked)
	bucket, err := s.Bucket("bucket")
	b := bucket.(*mock.Bucket)
	require.NoError(t, err)
	require.True(t, b.SaveInvoked)

	item, err := b.Get("key")
	require.NoError(t, err)
	require.Equal(t, []byte("data"), item.Data)
}
