package boltdb_test

import (
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	pathStore, cleanupStore := preparePath(t, "store.db")
	defer cleanupStore()

	s, err := boltdb.NewStore(pathStore)
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
}
