package fs_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/asdine/brazier/store/fs"
	"github.com/stretchr/testify/require"
)

func prepareDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	return dir, func() {
		os.RemoveAll(dir)
	}
}

func TestStore(t *testing.T) {
	dir, cleanup := prepareDir(t)
	defer cleanup()

	s := fs.NewStore(dir)

	require.Equal(t, "fs", s.Name())

	err := s.Create("bucket1")
	require.NoError(t, err)

	bucket, err := s.Bucket("bucket1")
	require.NoError(t, err)
	require.NotNil(t, bucket)

	err = bucket.Close()
	require.NoError(t, err)
}
