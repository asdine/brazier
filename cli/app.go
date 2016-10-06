package cli

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/config"
	"github.com/asdine/brazier/rpc/proto"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultDBName     = "brazier.db"
	defaultDataDir    = ".brazier"
	defaultSocketName = "brazier.sock"
	defaultBucket     = "default"
	settingsDB        = "settings.db"
	registryDB        = "registry.db"
)

// App is the main cli application
type app struct {
	Out        io.Writer
	Cli        Cli
	Registry   brazier.Registry
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

	if a.Store == nil {
		a.Store, err = boltdb.NewStore(filepath.Join(a.DataDir, defaultDBName))
		if err != nil {
			return err
		}
	}

	if a.Registry == nil {
		a.Registry, err = boltdb.NewRegistry(filepath.Join(a.DataDir, registryDB), a.Store)
		if err != nil {
			return err
		}
	}

	a.Cli = &cli{App: a}

	return nil
}

func (a *app) PostRun(cmd *cobra.Command, args []string) error {
	var err error

	if a.conn != nil {
		err = a.conn.Close()
		if err != nil {
			return err
		}
	}

	if a.Registry != nil {
		err = a.Registry.Close()
		if err != nil {
			return err
		}
	}

	if a.Store != nil {
		err = a.Store.Close()
		if err != nil {
			return err
		}
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

func (a *app) settingsDB() (*storm.DB, error) {
	return storm.Open(filepath.Join(a.DataDir, settingsDB))
}

func (a *app) defaultBucket() (string, error) {
	db, err := a.settingsDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var name string
	err = db.Get("buckets", "default", &name)
	if err != nil && err != storm.ErrNotFound {
		return "", err
	}

	if name == "" {
		name = defaultBucket
	}

	return name, nil
}
