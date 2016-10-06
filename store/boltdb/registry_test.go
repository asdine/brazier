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
	pathReg, cleanupReg := preparePath(t, "reg.db")
	defer cleanupReg()

	s, err := boltdb.NewStore(pathStore)
	require.NoError(t, err)

	r, err := boltdb.NewRegistry(pathReg, s)
	require.NoError(t, err)

	err = r.Create("bucket1")
	require.NoError(t, err)

	err = r.Create("bucket1")
	require.Error(t, err)
	require.Equal(t, store.ErrAlreadyExists, err)

	info1, err := r.BucketInfo("bucket1")
	require.NoError(t, err)
	require.NotNil(t, info1)
	require.Equal(t, "bucket1", info1.Name)

	b1, err := r.Bucket(info1.Name)
	require.NoError(t, err)
	require.NotNil(t, b1)
	defer b1.Close()

	_, err = r.BucketInfo("bucket2")
	require.Equal(t, err, store.ErrNotFound)

	err = r.Create("bucket2")
	require.NoError(t, err)

	info2, err := r.BucketInfo("bucket2")
	require.NoError(t, err)
	require.Equal(t, "bucket2", info2.Name)

	list, err := r.List()
	require.NoError(t, err)
	require.Len(t, list, 2)

	err = r.Close()
	require.NoError(t, err)
}
