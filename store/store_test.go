package store_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
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

	t.Run("EndingWithSlash", func(t *testing.T) {
		path, key := store.SplitPathKey("a/b/c/")
		require.Empty(t, key)
		require.Len(t, path, 3)
		require.Equal(t, "a", path[0])
		require.Equal(t, "b", path[1])
		require.Equal(t, "c", path[2])
	})
}

func TestStoreWithMock(t *testing.T) {
	testStore(t, "mock")
}

func TestStoreWithBoltDB(t *testing.T) {
	testStore(t, "boltdb")
}

func testStore(t *testing.T, backendType string) {
	t.Run("EmptyRegistry", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		_, err := s.Get("/")
		require.Equal(t, store.ErrForbidden, err)

		_, err = s.Get("/a/b")
		require.Equal(t, store.ErrNotFound, err)
	})

	t.Run("CreateBucket", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		err := s.CreateBucket("/")
		require.Equal(t, store.ErrAlreadyExists, err)

		err = s.CreateBucket("/a/b/c")
		require.Equal(t, store.ErrForbidden, err)

		err = s.CreateBucket("/a/b/c/")
		require.NoError(t, err)

		err = s.CreateBucket("a/b/c/")
		require.Equal(t, store.ErrAlreadyExists, err)
	})

	t.Run("Put", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		_, err := s.Put("/", []byte("Value"))
		require.Equal(t, store.ErrForbidden, err)

		item, err := s.Put("/a", []byte("Value"))
		require.NoError(t, err)
		require.NotNil(t, item)

		item, err = s.Put("/1a/2a", []byte("Value"))
		require.NoError(t, err)
		require.NotNil(t, item)

		item, err = s.Get("1a/2a")
		require.NoError(t, err)
		require.Equal(t, "2a", item.Key)
		require.Equal(t, []byte("Value"), item.Data)

		_, err = s.Put("/1a", []byte("Value"))
		require.NoError(t, err)

		item, err = s.Get("1a")
		require.NoError(t, err)
		require.Equal(t, "1a", item.Key)
		require.Equal(t, []byte("Value"), item.Data)
	})

	t.Run("Get", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		item, err := s.Put("/a/b/c", []byte("Value"))
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
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		for i := 0; i < 10; i++ {
			item, err := s.Put(fmt.Sprintf("/a/b/k%d", i), []byte("Value"+strconv.Itoa(i)))
			require.NoError(t, err)
			require.NotNil(t, item)
		}

		items, err := s.List("/a/c/", 1, 10)
		require.Equal(t, store.ErrNotFound, err)

		items, err = s.List("/a/b/", 1, 10)
		require.NoError(t, err)
		require.Len(t, items, 10)

		items, err = s.List("/a/b", 1, 10)
		require.Equal(t, store.ErrForbidden, err)
	})

	t.Run("Tree", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		for i := 0; i < 3; i++ {
			for j := 0; j < 5; j++ {
				item, err := s.Put(fmt.Sprintf("/a/b%d/k%d", i, j), []byte("Value"+strconv.Itoa(j)))
				require.NoError(t, err)
				require.NotNil(t, item)
			}
		}

		items, err := s.Tree("/")
		require.NoError(t, err)

		items, err = s.Tree("/z/")
		require.Equal(t, store.ErrNotFound, err)

		items, err = s.Tree("/a/c/")
		require.Equal(t, store.ErrNotFound, err)

		items, err = s.Tree("/a/")
		require.NoError(t, err)
		require.Len(t, items, 3)
		for i := 0; i < 3; i++ {
			require.Len(t, items[i].Children, 5)
			require.Equal(t, fmt.Sprintf("b%d/", i), items[i].Key)
			for j := 0; j < 5; j++ {
				require.Equal(t, []byte("Value"+strconv.Itoa(j)), items[i].Children[j].Data)
			}
		}

		items, err = s.Tree("/a")
		require.Equal(t, store.ErrForbidden, err)
	})

	t.Run("Delete", func(t *testing.T) {
		r, cleanup := getRegistryHelper(t, backendType)
		defer cleanup()
		s := store.NewStore(r)

		for i := 0; i < 10; i++ {
			item, err := s.Put(fmt.Sprintf("/a/b/k%d", i), []byte("Value"+strconv.Itoa(i)))
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

		// make sure we can't delete a bucket with Delete
		err = s.CreateBucket("/v/d/")
		require.NoError(t, err)

		err = s.Delete("/v/d/")
		require.Equal(t, store.ErrForbidden, err)
	})
}

func boltRegistryHelper(t *testing.T) (brazier.Registry, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	b, err := boltdb.NewBackend(path.Join(dir, "backend.db"))
	require.NoError(t, err)

	r, err := boltdb.NewRegistry(path.Join(dir, "registry.db"), b)
	require.NoError(t, err)

	return r, func() {
		b.Close()
		r.Close()
		os.RemoveAll(dir)
	}
}

func mockRegistryHelper(t *testing.T) brazier.Registry {
	return mock.NewRegistry(mock.NewBackend())
}

func getRegistryHelper(t *testing.T, backendType string) (brazier.Registry, func()) {
	switch backendType {
	case "boltdb":
		return boltRegistryHelper(t)
	default:
		return mockRegistryHelper(t), func() {}
	}
}
