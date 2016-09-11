package boltdb_test

import (
	"fmt"
	"testing"

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

	i, err := b.Save("id", []byte("Data"))
	require.NoError(t, err)
	require.Equal(t, "id", i.Key)
	require.Equal(t, []byte("Data"), i.Data)

	j, err := b.Save("id", []byte("New Data"))
	require.NoError(t, err)
	require.Equal(t, []byte("New Data"), j.Data)

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

	j, err := b.Get(i.Key)
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

	_, err = b.Get(i.Key)
	require.NoError(t, err)

	err = b.Delete(i.Key)
	require.NoError(t, err)

	err = b.Delete(i.Key)
	require.NoError(t, err)

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
		_, err := b.Save(fmt.Sprintf("%c", i+65), []byte("Data"))
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
	require.Equal(t, "A", list[0].Key)
	require.Equal(t, "E", list[4].Key)

	list, err = b.Page(1, 25)
	require.NoError(t, err)
	require.Len(t, list, 20)
	require.Equal(t, "A", list[0].Key)
	require.Equal(t, "T", list[19].Key)

	list, err = b.Page(2, 5)
	require.NoError(t, err)
	require.Len(t, list, 5)
	require.Equal(t, "F", list[0].Key)
	require.Equal(t, "J", list[4].Key)

	list, err = b.Page(2, 15)
	require.NoError(t, err)
	require.Len(t, list, 5)
	require.Equal(t, "P", list[0].Key)
	require.Equal(t, "T", list[4].Key)

	list, err = b.Page(3, 15)
	require.NoError(t, err)
	require.Len(t, list, 0)

	// all
	list, err = b.Page(1, -1)
	require.NoError(t, err)
	require.Len(t, list, 20)
	require.Equal(t, "A", list[0].Key)
	require.Equal(t, "T", list[19].Key)
}
