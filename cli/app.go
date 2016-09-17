package cli

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/config"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/spf13/cobra"
)

const (
	defaultDBName     = "brazier.db"
	defaultDataDir    = ".brazier"
	defaultSocketName = "brazier.sock"
)

// App is the main cli application
type app struct {
	Out        io.Writer
	Cli        Cli
	Store      brazier.Store
	ConfigPath string
	DataDir    string
	SocketPath string
	Config     config.Config
	conn       *grpc.ClientConn
}

// Run runs the root command
func (a *app) Run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// PreRun runs the root command
func (a *app) PreRun(cmd *cobra.Command, args []string) error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	err = a.initDataDir()
	if err != nil {
		return err
	}

	a.SocketPath = filepath.Join(a.DataDir, defaultSocketName)

	if a.serverIsLaunched() {
		client, err := a.rpcClient()
		if err != nil {
			return err
		}

		a.Cli = &rpcCli{
			App:    a,
			Client: client,
		}
		return nil
	}

	a.Store, err = boltdb.NewStore(filepath.Join(a.DataDir, defaultDBName))
	if err != nil {
		return err
	}
	a.Cli = &cli{App: a}

	return nil
}

func (a *app) PostRun(cmd *cobra.Command, args []string) error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// manages brazier config
func (a *app) initConfig() error {
	if a.ConfigPath != "" {
		return config.FromFile(a.ConfigPath, &a.Config)
	}
	return nil
}

// manages brazier config
func (a *app) initDataDir() error {
	if a.DataDir == "" {
		// check in the local directory
		fi, err := os.Stat(defaultDataDir)
		if err == nil && fi.Mode().IsDir() {
			a.DataDir = defaultDataDir
			return nil
		}

		// check in the home directory
		home := os.Getenv("HOME")
		if home == "" {
			return errors.New("Can't find $HOME directory")
		}
		a.DataDir = filepath.Join(home, defaultDataDir)
	}

	fi, err := os.Stat(a.DataDir)
	if err != nil {
		return os.Mkdir(a.DataDir, 0755)
	}

	if !fi.Mode().IsDir() {
		return errors.New("Data directory must be a valid directory")
	}

	return nil
}

func (a *app) serverIsLaunched() bool {
	_, err := os.Stat(a.SocketPath)
	return err == nil
}

func (a *app) rpcClient() (proto.BucketClient, error) {
	conn, err := grpc.Dial("",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			sock, err := net.DialTimeout("unix", a.SocketPath, timeout)
			return sock, err
		}),
	)
	if err != nil {
		return nil, err
	}

	a.conn = conn
	return proto.NewBucketClient(conn), nil
}
