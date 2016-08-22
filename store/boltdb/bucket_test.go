package boltdb_test

import (
	"testing"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/stretchr/testify/require"
)

func TestBucketAdd(t *testing.T) {
	db, cleanup := prepareDB(t, storm.AutoIncrement())
	defer cleanup()

	var b brazier.Bucket
	node := db.From("buckets")
	b = boltdb.NewBucket(node.From("b1"))

	now := time.Now()
	i, err := b.Add([]byte("Data"), "json", "")
	require.NoError(t, err)
	require.True(t, i.CreatedAt.After(now))
	require.NotEmpty(t, i.ID)
	require.Equal(t, "json", i.MimeType)
	require.Equal(t, []byte("Data"), i.Data)

	i, err = b.Add([]byte("Data"), "json", "name")
	require.NoError(t, err)
	require.Equal(t, "name", i.ID)
}

func TestBucketGet(t *testing.T) {
	db, cleanup := prepareDB(t, storm.AutoIncrement())
	defer cleanup()

	node := db.From("buckets")
	b := boltdb.NewBucket(node.From("b1"))

	i, err := b.Add([]byte("Data"), "json", "")
	require.NoError(t, err)

	j, err := b.Get(i.ID)
	require.NoError(t, err)
	require.Equal(t, i, j)

	_, err = b.Get("some id")
	require.Equal(t, store.ErrNotFound, err)
}

func TestBucketDelete(t *testing.T) {
	db, cleanup := prepareDB(t, storm.AutoIncrement())
	defer cleanup()

	node := db.From("buckets")
	b := boltdb.NewBucket(node.From("b1"))

	i, err := b.Add([]byte("Data"), "json", "")
	require.NoError(t, err)

	_, err = b.Get(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.Equal(t, store.ErrNotFound, err)
}
