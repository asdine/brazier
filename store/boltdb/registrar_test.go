package boltdb_test

import (
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/stretchr/testify/require"
)

func TestRegistrar(t *testing.T) {
	db, cleanup := prepareDB(t)
	defer cleanup()

	r := boltdb.NewRegistrar(db)
	s := boltdb.NewStore(db)

	info1, err := r.Create("b1", s)
	require.NoError(t, err)

	info2, err := r.Create(info1.ID, s)
	require.Error(t, err)
	require.Equal(t, store.ErrAlreadyExists, err)

	err = s.Create(info1.ID)
	require.NoError(t, err)

	info2, err = r.Bucket(info1.ID)
	require.NoError(t, err)
	require.Equal(t, info1, info2)

	_, err = r.Bucket("something")
	require.Error(t, err)
	require.Equal(t, store.ErrNotFound, err)
}
