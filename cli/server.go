package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/http"
	"github.com/asdine/brazier/rpc"
	"github.com/spf13/cobra"
)

// NewServerCmd creates a "Server" cli command
func NewServerCmd(a *app) *cobra.Command {
	serverCmd := serverCmd{
		App:              a,
		HTTPServerFunc:   http.NewServer,
		RPCServerFunc:    rpc.NewServer,
		SocketServerFunc: rpc.NewServer,
		useExit:          true,
		c:                make(chan os.Signal, 1),
	}

	cmd := cobra.Command{
		Use:   "server",
		Short: "Run Brazier as an HTTP and RPC server",
		Long:  "Run Brazier as an HTTP and RPC server",
		RunE:  serverCmd.Serve,
	}

	cmd.Flags().StringVar(&serverCmd.App.Config.HTTP.Address, "http-addr", ":5656", "HTTP address")
	cmd.Flags().StringVar(&serverCmd.App.Config.RPC.Address, "rpc-addr", "127.0.0.1:5657", "RPC address")
	return &cmd
}

type serverCmd struct {
	App              *app
	c                chan os.Signal
	useExit          bool
	HTTPServerFunc   func(brazier.Store) brazier.Server
	RPCServerFunc    func(brazier.Store) brazier.Server
	SocketServerFunc func(brazier.Store) brazier.Server
}

func (s *serverCmd) Serve(cmd *cobra.Command, args []string) error {
	servers, err := s.createServers()
	if err != nil {
		return err
	}
	fmt.Fprintf(s.App.Out, "Serving HTTP on address %s\n", s.App.Config.HTTP.Address)
	fmt.Fprintf(s.App.Out, "Serving RPC on address %s\n", s.App.Config.RPC.Address)

	s.runServers(servers)
	return nil
}

func (s *serverCmd) createServers() (map[net.Listener]brazier.Server, error) {
	servers := make(map[net.Listener]brazier.Server)

	httpListener, err := net.Listen("tcp", s.App.Config.HTTP.Address)
	if err != nil {
		return nil, err
	}
	servers[httpListener] = s.HTTPServerFunc(s.App.Store)

	rpcListener, err := net.Listen("tcp", s.App.Config.RPC.Address)
	if err != nil {
		return nil, err
	}
	servers[rpcListener] = s.RPCServerFunc(s.App.Store)

	socketListener, err := net.Listen("unix", filepath.Join(s.App.DataDir, "brazier.sock"))
	if err != nil {
		return nil, err
	}
	servers[socketListener] = s.SocketServerFunc(s.App.Store)

	return servers, nil
}

func (s *serverCmd) runServers(servers map[net.Listener]brazier.Server) {
	var wg sync.WaitGroup

	for l, srv := range servers {
		wg.Add(1)
		go func(l net.Listener, srv brazier.Server) {
			defer wg.Done()
			srv.Serve(l)
		}(l, srv)
	}

	signal.Notify(s.c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range s.c {
			fmt.Fprintf(s.App.Out, "\nStopping servers...")
			for _, srv := range servers {
				srv.Stop(time.Second)
			}
			fmt.Fprintf(s.App.Out, " OK\n")
			if s.useExit {
				os.Exit(1)
			}
		}
	}()

	wg.Wait()
}
