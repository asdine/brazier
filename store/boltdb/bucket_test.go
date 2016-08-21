package boltdb_test

import (
	"testing"
	"time"

	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/stretchr/testify/require"
)

func TestBucketAdd(t *testing.T) {
	db, cleanup := prepareDB(t, storm.AutoIncrement())
	defer cleanup()

	node := db.From("buckets")
	b := boltdb.NewBucket(node.From("b1"))

	now := time.Now()
	i, err := b.Add([]byte("Data"), "json", "")
	require.NoError(t, err)
	require.True(t, i.CreatedAt.After(now))
	require.Zero(t, i.UpdatedAt)
	require.NotEmpty(t, i.ID)
	require.Equal(t, "json", i.MimeType)
	require.Equal(t, []byte("Data"), i.Data)

	i, err = b.Add([]byte("Data"), "json", "name")
	require.NoError(t, err)
	require.Equal(t, "name", i.ID)
}
