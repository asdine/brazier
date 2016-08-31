package boltdb_test

import (
	"testing"
	"time"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/stretchr/testify/require"
)

func TestBucketSave(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	s := boltdb.NewStore(path)

	b, err := s.Bucket("b1")
	require.NoError(t, err)

	now := time.Now()
	i, err := b.Save("id", []byte("Data"))
	require.NoError(t, err)
	require.True(t, i.CreatedAt.After(now))
	require.Equal(t, "id", i.ID)
	require.Equal(t, []byte("Data"), i.Data)

	j, err := b.Save("id", []byte("New Data"))
	require.NoError(t, err)
	require.Equal(t, i.CreatedAt, j.CreatedAt)
	require.Equal(t, []byte("New Data"), j.Data)
	require.True(t, j.UpdatedAt.After(j.CreatedAt))

	err = b.Close()
	require.NoError(t, err)
}

func TestBucketGet(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	s := boltdb.NewStore(path)

	b, err := s.Bucket("b1")
	require.NoError(t, err)

	i, err := b.Save("id", []byte("Data"))
	require.NoError(t, err)

	j, err := b.Get(i.ID)
	require.NoError(t, err)
	require.Equal(t, i.Data, j.Data)

	_, err = b.Get("some id")
	require.Equal(t, store.ErrNotFound, err)

	err = b.Close()
	require.NoError(t, err)
}

func TestBucketDelete(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	s := boltdb.NewStore(path)

	b, err := s.Bucket("b1")
	require.NoError(t, err)

	i, err := b.Save("id", []byte("Data"))
	require.NoError(t, err)

	_, err = b.Get(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.Equal(t, store.ErrNotFound, err)

	err = b.Close()
	require.NoError(t, err)
}
