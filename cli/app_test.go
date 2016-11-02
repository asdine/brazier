package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/rpc"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func testableApp(t *testing.T) (*app, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	require.NoError(t, err)

	a := app{
		Out:        bytes.NewBuffer([]byte("")),
		Store:      store.NewStore(mock.NewRegistry(mock.NewBackend())),
		DataDir:    dir,
		SocketPath: filepath.Join(dir, defaultSocketName),
	}

	a.PreRun(nil, nil)

	return &a, func() {
		a.PostRun(nil, nil)
		os.RemoveAll(dir)
	}
}

func testableAppRPC(t *testing.T) (*app, func()) {
	app, cleanup := testableApp(t)

	app.Config.HTTP.Address = ":"
	app.Config.RPC.Address = ":"

	s := serverCmd{
		App:              app,
		HTTPServerFunc:   mock.NewServer,
		RPCServerFunc:    mock.NewServer,
		SocketServerFunc: rpc.NewServer,
		c:                make(chan os.Signal, 1),
	}

	servers, err := s.createServers()
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runServers(servers)
	}()

	app.PreRun(nil, nil)

	return app, func() {
		cleanup()
		s.c <- os.Interrupt
		wg.Wait()
	}
}

func TestAppDataDir(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	dir := app.DataDir
	app.DataDir = ""

	// using HOME directory
	os.Setenv("HOME", dir)
	err := app.initDataDir()
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
