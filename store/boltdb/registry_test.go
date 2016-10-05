package boltdb_test

import (
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	r, err := boltdb.NewRegistry(path)
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
