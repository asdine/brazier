package boltdb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/stretchr/testify/require"
)

type errorHandler interface {
	Error(args ...interface{})
}

func preparePath(t errorHandler, dbName string) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	if err != nil {
		t.Error(err)
	}

	return filepath.Join(dir, dbName), func() {
		os.RemoveAll(dir)
	}
}

func prepareDB(t *testing.T, dbName string, opts ...func(*storm.DB) error) (*storm.DB, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	db, err := storm.Open(filepath.Join(dir, dbName), opts...)
	require.NoError(t, err)

	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func TestBackend(t *testing.T) {
	path, cleanup := preparePath(t, "backend.db")
	defer cleanup()

	s, err := boltdb.NewBackend(path)
	require.NoError(t, err)
	defer s.Close()

	bucket, err := s.Bucket("a")
	require.NoError(t, err)
	require.NotNil(t, bucket)
	require.NotNil(t, s.DB)

	err = bucket.Close()
	require.NoError(t, err)

	b1, err := s.Bucket("a")
	require.NoError(t, err)

	b2, err := s.Bucket("b")
	require.NoError(t, err)

	b1bis, err := s.Bucket("a", "b", "c")
	require.NoError(t, err)

	err = b1.Close()
	require.NoError(t, err)

	err = b2.Close()
	require.NoError(t, err)

	err = b1bis.Close()
	require.NoError(t, err)
}
