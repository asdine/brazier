package boltdb_test

import (
	"bytes"
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/boltdb/bolt"
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

		r.DB.Bolt.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("a"))
			require.NotNil(t, b)
			c := b.Cursor()
			var notEmpty bool
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if bytes.HasPrefix(k, []byte("__storm")) {
					continue
				}
				notEmpty = true
				break
			}

			require.False(t, notEmpty)
			return nil
		})

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

		b, err := r.Bucket("a")
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

		_, err = r.Children()
		require.Equal(t, store.ErrNotFound, err)

		err = r.Create("1a", "2a", "3a")
		require.NoError(t, err)

		err = r.Create("1a", "2a", "3b")
		require.NoError(t, err)

		tree, err := r.Children("1a")
		require.NoError(t, err)
		require.Len(t, tree, 1)
		require.Len(t, tree[0].Children, 2)
		require.Equal(t, "3a", tree[0].Children[0].Key)
		require.Equal(t, "3b", tree[0].Children[1].Key)

		tree, err = r.Children("1a", "2a")
		require.NoError(t, err)
		require.Len(t, tree, 2)
		require.Equal(t, "3a", tree[0].Key)
		require.Equal(t, "3b", tree[1].Key)

		tree, err = r.Children("1a", "2a", "3a")
		require.NoError(t, err)
		require.Len(t, tree, 0)

		_, err = r.Children("1a", "2a", "3c")
		require.Equal(t, store.ErrNotFound, err)
	})
}
