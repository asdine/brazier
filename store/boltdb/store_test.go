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

type errorHandler interface {
	Error(args ...interface{})
}

func preparePath(t errorHandler) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	if err != nil {
		t.Error(err)
	}

	return filepath.Join(dir, "brazier.db"), func() {
		os.RemoveAll(dir)
	}
}

func prepareDB(t *testing.T, opts ...func(*storm.DB) error) (*storm.DB, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	db, err := storm.Open(filepath.Join(dir, "brazier.db"), opts...)
	require.NoError(t, err)

	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func TestStore(t *testing.T) {
	path, cleanup := preparePath(t)
	defer cleanup()

	s, err := boltdb.NewStore(path)
	require.NoError(t, err)

	err = s.Create("bucket1")
	require.NoError(t, err)

	err = s.Create("bucket1")
	require.Error(t, err)
	require.Equal(t, store.ErrAlreadyExists, err)

	bucket, err := s.Bucket("bucket1")
	require.NoError(t, err)
	require.NotNil(t, bucket)
	require.NotNil(t, s.DB)

	err = bucket.Close()
	require.NoError(t, err)

	b1, err := s.Bucket("bucket1")
	require.NoError(t, err)

	b2, err := s.Bucket("bucket2")
	require.Equal(t, err, store.ErrNotFound)

	err = s.Create("bucket2")
	require.NoError(t, err)

	b2, err = s.Bucket("bucket2")
	require.NoError(t, err)

	b1bis, err := s.Bucket("bucket1")
	require.NoError(t, err)

	err = b1.Close()
	require.NoError(t, err)

	err = b2.Close()
	require.NoError(t, err)

	err = b1bis.Close()
	require.NoError(t, err)

	list, err := s.List()
	require.NoError(t, err)
	require.Len(t, list, 2)
}
