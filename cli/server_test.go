package cli

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/asdine/brazier/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateServers(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	app.Config.HTTP.Address = ":"
	app.Config.RPC.Address = ":"

	s := serverCmd{
		App:              app,
		HTTPServerFunc:   mock.NewServer,
		RPCServerFunc:    mock.NewServer,
		SocketServerFunc: mock.NewServer,
	}

	servers, err := s.createServers()
	require.NoError(t, err)
	require.Len(t, servers, 3)

	for l := range servers {
		err = l.Close()
		require.NoError(t, err)
	}
}

func TestRunServers(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	app.Config.HTTP.Address = ":"
	app.Config.RPC.Address = ":"

	s := serverCmd{
		App:              app,
		HTTPServerFunc:   mock.NewServer,
		RPCServerFunc:    mock.NewServer,
		SocketServerFunc: mock.NewServer,
		c:                make(chan os.Signal, 1),
	}

	servers, err := s.createServers()
	require.NoError(t, err)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runServers(servers)
		for _, srv := range servers {
			m := srv.(*mock.Server)
			require.True(t, m.ServeInvoked)
			require.True(t, m.StopInvoked)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	s.c <- os.Interrupt
	wg.Wait()
}
