package store_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func TestSplitPathKey(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		path, key := store.SplitPathKey("")
		require.Empty(t, path)
		require.Zero(t, key)
	})

	t.Run("Slash", func(t *testing.T) {
		path, key := store.SplitPathKey("/")
		require.Empty(t, path)
		require.Zero(t, key)
	})

	t.Run("Spaces", func(t *testing.T) {
		path, key := store.SplitPathKey(" aaaa bbbb cccc")
		require.Empty(t, path)
		require.Equal(t, key, " aaaa bbbb cccc")
	})

	t.Run("PathWithSlash", func(t *testing.T) {
		path, key := store.SplitPathKey("/a/b/c")
		require.Len(t, path, 2)
		require.Equal(t, "a", path[0])
		require.Equal(t, "b", path[1])
		require.Equal(t, key, "c")
	})

	t.Run("PathWithoutSlash", func(t *testing.T) {
		path, key := store.SplitPathKey("a/b/c")
		require.Len(t, path, 2)
		require.Equal(t, "a", path[0])
		require.Equal(t, "b", path[1])
		require.Equal(t, key, "c")
	})
}

func TestStore(t *testing.T) {
	t.Run("EmptyRegistry", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		_, err := s.Get("/")
		require.Equal(t, store.ErrNotFound, err)

		_, err = s.Get("/a/b")
		require.Equal(t, store.ErrNotFound, err)
	})

	t.Run("CreateBucket", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		err := s.CreateBucket("/")
		require.Equal(t, store.ErrAlreadyExists, err)

		err = s.CreateBucket("/a/b/c")
		require.NoError(t, err)
	})

	t.Run("Save", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		_, err := s.Save("/", []byte("Value"))
		require.Equal(t, store.ErrForbidden, err)

		_, err = s.Save("/a", []byte("Value"))
		require.Equal(t, store.ErrForbidden, err)

		item, err := s.Save("/a/b/c", []byte("Value"))
		require.NoError(t, err)
		require.NotNil(t, item)

		// TODO enforce the following behaviour
		// _, err = s.Save("/a/b", []byte("Value"))
		// require.Equal(t, store.ErrAlreadyExists, err)
	})

	t.Run("Get", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		item, err := s.Save("/a/b/c", []byte("Value"))
		require.NoError(t, err)
		require.NotNil(t, item)

		found, err := s.Get("/a/b/c")
		require.NoError(t, err)
		require.Equal(t, item, found)

		found, err = s.Get("a/b/c")
		require.NoError(t, err)
		require.Equal(t, item, found)
	})

	t.Run("List", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		for i := 0; i < 10; i++ {
			item, err := s.Save(fmt.Sprintf("/a/b/k%d", i), []byte("Value"+strconv.Itoa(i)))
			require.NoError(t, err)
			require.NotNil(t, item)
		}

		items, err := s.List("/a/c", 1, 10)
		require.Equal(t, store.ErrNotFound, err)

		items, err = s.List("/a/b", 1, 10)
		require.NoError(t, err)
		require.Len(t, items, 10)
	})

	t.Run("Delete", func(t *testing.T) {
		r := mock.NewRegistry(mock.NewBackend())
		s := store.NewStore(r)

		for i := 0; i < 10; i++ {
			item, err := s.Save(fmt.Sprintf("/a/b/k%d", i), []byte("Value"+strconv.Itoa(i)))
			require.NoError(t, err)
			require.NotNil(t, item)
		}

		err := s.Delete("/a/b")
		require.Equal(t, store.ErrNotFound, err)

		err = s.Delete("/a/b/k10")
		require.Equal(t, store.ErrNotFound, err)

		err = s.Delete("/a/b/k5")
		require.NoError(t, err)

		err = s.Delete("/a/b/k5")
		require.Equal(t, store.ErrNotFound, err)
	})
}
