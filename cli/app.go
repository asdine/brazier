package cli

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/config"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New() *cobra.Command {
	a := app{Out: os.Stdout}

	cmd := cobra.Command{
		Use:               "brazier",
		Short:             "Brazier",
		Long:              `Brazier`,
		Run:               a.Run,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: a.PreRun,
	}

	cmd.SetOutput(os.Stdout)
	cmd.AddCommand(NewCreateCmd(&a))
	cmd.AddCommand(NewSaveCmd(&a))
	cmd.AddCommand(NewGetCmd(&a))
	cmd.AddCommand(NewDeleteCmd(&a))
	cmd.AddCommand(NewListCmd(&a))
	cmd.AddCommand(NewHTTPCmd(&a))
	cmd.AddCommand(NewRPCCmd(&a))

	cmd.PersistentFlags().StringVarP(&a.ConfigPath, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVar(&a.DataDir, "data-dir", "", "data directory (default $HOME/.brazier)")
	return &cmd
}

// App is the main cli application
type app struct {
	Out        io.Writer
	Store      brazier.Store
	ConfigPath string
	DataDir    string
	Config     config.Config
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

	a.Store = boltdb.NewStore(filepath.Join(a.DataDir, "brazier.db"))
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
		fi, err := os.Stat(".brazier")
		if err == nil && fi.Mode().IsDir() {
			a.DataDir = ".brazier"
			return nil
		}

		// check in the home directory
		home := os.Getenv("HOME")
		if home == "" {
			return errors.New("Can't find $HOME directory")
		}
		a.DataDir = filepath.Join(home, ".brazier")
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
