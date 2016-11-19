package boltdb_test

import (
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	pathStore, cleanupStore := preparePath(t, "backend.db")
	defer cleanupStore()

	s, err := boltdb.NewBackend(pathStore)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create", func(t *testing.T) {
		pathReg, cleanupReg := preparePath(t, "reg.db")
		defer cleanupReg()
		r, err := boltdb.NewRegistry(pathReg, s)
		require.NoError(t, err)
		defer r.Close()

		err = r.Create("a")
		require.NoError(t, err)

		err = r.Create("a")
		require.Equal(t, store.ErrAlreadyExists, err)

		err = r.Create("a", "b")
		require.NoError(t, err)

		err = r.Create("a", "b")
		require.Equal(t, store.ErrAlreadyExists, err)

		err = r.Create("e", "f", "g", "h")
		require.NoError(t, err)

		err = r.Create("e", "f")
		require.Equal(t, store.ErrAlreadyExists, err)
	})

	t.Run("bucket", func(t *testing.T) {
		pathReg, cleanupReg := preparePath(t, "reg.db")
		defer cleanupReg()
		r, err := boltdb.NewRegistry(pathReg, s)
		require.NoError(t, err)
		defer r.Close()

		b, err := r.Bucket()
		require.NoError(t, err)

		b, err = r.Bucket("a")
		require.Equal(t, store.ErrNotFound, err)

		b, err = r.Bucket("a", "b")
		require.Equal(t, store.ErrNotFound, err)

		err = r.Create("a")
		require.NoError(t, err)

		b, err = r.Bucket("a")
		require.NoError(t, err)
		require.NotNil(t, b)

		b, err = r.Bucket("a", "b")
		require.Equal(t, store.ErrNotFound, err)

		err = r.Create("a", "b")
		require.NoError(t, err)

		b, err = r.Bucket("a", "b")
		require.NoError(t, err)
		require.NotNil(t, b)
	})

	t.Run("tree", func(t *testing.T) {
		pathReg, cleanupReg := preparePath(t, "reg.db")
		defer cleanupReg()
		r, err := boltdb.NewRegistry(pathReg, s)
		require.NoError(t, err)
		defer r.Close()

		_, err = r.Children("a", "b")
		require.Equal(t, store.ErrNotFound, err)

		tree, err := r.Children()
		require.NoError(t, err)
		require.Len(t, tree, 0)

		err = r.Create("1a", "2a", "3b")
		require.NoError(t, err)

		err = r.Create("1a", "2a", "3a")
		require.NoError(t, err)

		tree, err = r.Children()
		require.NoError(t, err)
		require.Len(t, tree, 1)

		tree, err = r.Children("1a")
		require.NoError(t, err)
		require.Len(t, tree, 1)
		require.Len(t, tree[0].Children, 2)
		require.Equal(t, "3b", tree[0].Children[0].Key)
		require.Equal(t, "3a", tree[0].Children[1].Key)

		tree, err = r.Children("1a", "2a")
		require.NoError(t, err)
		require.Len(t, tree, 2)
		// check if the insert order is preserved
		require.Equal(t, "3b", tree[0].Key)
		require.Equal(t, "3a", tree[1].Key)

		tree, err = r.Children("1a", "2a", "3a")
		require.NoError(t, err)
		require.Len(t, tree, 0)

		_, err = r.Children("1a", "2a", "3c")
		require.Equal(t, store.ErrNotFound, err)

		// all children from root
		tree, err = r.Children()
		require.NoError(t, err)
		require.Len(t, tree, 1)
		require.Len(t, tree[0].Children, 1)
		require.Len(t, tree[0].Children[0].Children, 2)
		require.Equal(t, "3b", tree[0].Children[0].Children[0].Key)
		require.Equal(t, "3a", tree[0].Children[0].Children[1].Key)
	})
}
