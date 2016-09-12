package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/stretchr/testify/require"
)

func testableApp(t *testing.T) *app {
	return &app{Out: bytes.NewBuffer([]byte("")), Store: mock.NewStore()}
}

func TestAppDataDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	app := testableApp(t)

	// using HOME directory
	os.Setenv("HOME", dir)
	err = app.initDataDir()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, ".brazier"), app.DataDir)
	fi, err := os.Stat(app.DataDir)
	require.NoError(t, err)
	require.True(t, fi.Mode().IsDir())
	require.Equal(t, os.FileMode(0755), fi.Mode().Perm())

	// already exists and valid
	err = app.initDataDir()
	require.NoError(t, err)

	// already exists and invalid
	err = os.Remove(app.DataDir)
	require.NoError(t, err)
	_, err = os.Create(app.DataDir)
	require.NoError(t, err)
	err = app.initDataDir()
	require.Error(t, err)

	// using specific directory
	app.DataDir = filepath.Join(dir, "data")
	err = app.initDataDir()
	require.NoError(t, err)
	fi, err = os.Stat(app.DataDir)
	require.NoError(t, err)
	require.True(t, fi.Mode().IsDir())
	require.Equal(t, os.FileMode(0755), fi.Mode().Perm())
}
