package fs_test

import (
	"strings"
	"testing"
	"time"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/fs"
	"github.com/stretchr/testify/require"
)

func TestBucketAdd(t *testing.T) {
	dir, cleanup := prepareDir(t)
	defer cleanup()

	s := fs.NewStore(dir)
	err := s.Create("b1")
	require.NoError(t, err)
	b, err := s.Bucket("b1")
	require.NoError(t, err)

	now := time.Now()
	i, err := b.Add([]byte("Data"), "application/json", "")
	require.NoError(t, err)
	require.NotNil(t, i)
	require.True(t, i.CreatedAt.After(now))
	require.True(t, strings.HasSuffix(i.ID, ".json"))
	require.Equal(t, "application/json", i.MimeType)
	require.Equal(t, []byte("Data"), i.Data)

	i, err = b.Add([]byte("Data"), "application/json", "name.json")
	require.NoError(t, err)
	require.Equal(t, "name.json", i.ID)

	i, err = b.Add([]byte("Data"), "something", "")
	require.Error(t, err)
}

func TestBucketGet(t *testing.T) {
	dir, cleanup := prepareDir(t)
	defer cleanup()

	s := fs.NewStore(dir)
	err := s.Create("b1")
	require.NoError(t, err)
	b, err := s.Bucket("b1")
	require.NoError(t, err)

	i, err := b.Add([]byte("Data"), "application/json", "")
	require.NoError(t, err)

	j, err := b.Get(i.ID)
	require.NoError(t, err)
	require.Equal(t, i, j)

	_, err = b.Get("some id")
	require.Equal(t, store.ErrNotFound, err)
}

func TestBucketDelete(t *testing.T) {
	dir, cleanup := prepareDir(t)
	defer cleanup()

	s := fs.NewStore(dir)
	err := s.Create("b1")
	require.NoError(t, err)
	b, err := s.Bucket("b1")
	require.NoError(t, err)

	i, err := b.Add([]byte("Data"), "application/json", "")
	require.NoError(t, err)

	_, err = b.Get(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.NoError(t, err)

	err = b.Delete(i.ID)
	require.Equal(t, store.ErrNotFound, err)
}
