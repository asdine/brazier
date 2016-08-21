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

	info, err := s.Create("bucket1")
	require.NoError(t, err)
	require.Equal(t, "bucket1", info.ID)
	require.Equal(t, s.Name(), info.Store)

	bucket, err := s.Bucket(info.ID)
	require.NoError(t, err)
	require.NotNil(t, bucket)

	err = bucket.Close()
	require.NoError(t, err)
}
