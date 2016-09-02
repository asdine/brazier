package boltdb_test

import (
	"fmt"
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

func TestBucketPage(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	s := boltdb.NewStore(path)

	b, err := s.Bucket("b1")
	require.NoError(t, err)
	defer b.Close()

	for i := 0; i < 20; i++ {
		_, err := b.Save(fmt.Sprintf("id%d", i), []byte("Data"))
		require.NoError(t, err)
	}

	list, err := b.Page(0, 0)
	require.NoError(t, err)
	require.Len(t, list, 0)

	list, err = b.Page(0, 10)
	require.NoError(t, err)
	require.Len(t, list, 0)

	list, err = b.Page(1, 5)
	require.NoError(t, err)
	require.Len(t, list, 5)
	require.Equal(t, "id0", list[0].ID)
	require.Equal(t, "id4", list[4].ID)

	list, err = b.Page(1, 25)
	require.NoError(t, err)
	require.Len(t, list, 20)
	require.Equal(t, "id0", list[0].ID)
	require.Equal(t, "id19", list[19].ID)

	list, err = b.Page(2, 5)
	require.NoError(t, err)
	require.Len(t, list, 5)
	require.Equal(t, "id5", list[0].ID)
	require.Equal(t, "id9", list[4].ID)

	list, err = b.Page(2, 15)
	require.NoError(t, err)
	require.Len(t, list, 5)
	require.Equal(t, "id15", list[0].ID)
	require.Equal(t, "id19", list[4].ID)

	list, err = b.Page(3, 15)
	require.NoError(t, err)
	require.Len(t, list, 0)

	// all
	list, err = b.Page(1, -1)
	require.NoError(t, err)
	require.Len(t, list, 20)
	require.Equal(t, "id0", list[0].ID)
	require.Equal(t, "id19", list[19].ID)
}
