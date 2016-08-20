package boltdb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/stretchr/testify/require"
)

func prepareDB(t *testing.T) (*storm.DB, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	db, err := storm.Open(filepath.Join(dir, "brazier.db"))
	require.NoError(t, err)

	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func TestStore(t *testing.T) {
	db, cleanup := prepareDB(t)
	defer cleanup()

	s := boltdb.NewStore(db)

	require.Equal(t, "boltdb", s.Name())

	info, err := s.Create("bucket1")
	require.NoError(t, err)
	require.Equal(t, "bucket1", info.ID)
	require.Equal(t, s.Name(), info.Store)

	_, err = s.Create("bucket1")
	require.Error(t, err)
	require.Equal(t, store.ErrAlreadyExists, err)

	bucket, err := s.Bucket(info.ID)
	require.NoError(t, err)

	err = bucket.Close()
	require.NoError(t, err)

	bucket, err = s.Bucket("some id")
	require.Error(t, err)
	require.Equal(t, store.ErrNotFound, err)
}
