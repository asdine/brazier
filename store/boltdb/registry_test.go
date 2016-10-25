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

		err = r.Create("bucket1")
		require.NoError(t, err)

		err = r.Create("bucket1")
		require.Equal(t, store.ErrAlreadyExists, err)

		err = r.Create("bucket1", "bucket2")
		require.NoError(t, err)

		err = r.Create("bucket1", "bucket2")
		require.Equal(t, store.ErrAlreadyExists, err)
	})

	t.Run("bucketConfig", func(t *testing.T) {
		pathReg, cleanupReg := preparePath(t, "reg.db")
		defer cleanupReg()
		r, err := boltdb.NewRegistry(pathReg, s)
		require.NoError(t, err)
		defer r.Close()

		err = r.Create("bucket1", "bucket2")
		require.NoError(t, err)

		info1, err := r.BucketConfig("bucket1", "bucket2")
		require.NoError(t, err)
		require.NotNil(t, info1)
		require.Equal(t, []string{"bucket1", "bucket2"}, info1.Path)

		b1, err := r.Bucket(info1.Path...)
		require.NoError(t, err)
		require.NotNil(t, b1)
		defer b1.Close()

		_, err = r.BucketConfig("bucket2")
		require.Equal(t, err, store.ErrNotFound)

		err = r.Create("bucket2")
		require.NoError(t, err)

		info2, err := r.BucketConfig("bucket2")
		require.NoError(t, err)
		require.Equal(t, []string{"bucket2"}, info2.Path)
	})

	t.Run("list", func(t *testing.T) {
		pathReg, cleanupReg := preparePath(t, "reg.db")
		defer cleanupReg()
		r, err := boltdb.NewRegistry(pathReg, s)
		require.NoError(t, err)
		defer r.Close()

		list, err := r.List()
		require.NoError(t, err)
		require.Len(t, list, 0)

		err = r.Create("bucket1")
		require.NoError(t, err)

		err = r.Create("bucket1", "bucket2")
		require.NoError(t, err)

		err = r.Create("bucket2", "bucket1")
		require.NoError(t, err)

		list, err = r.List()
		require.NoError(t, err)
		require.Len(t, list, 3)
	})
}
