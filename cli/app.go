package cli

import (
	"io"
	"os"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/asdine/storm"
	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New() *cobra.Command {
	a := app{Out: os.Stdout, Path: "brazier.db"}

	cmd := cobra.Command{
		Use:                "brazier",
		Short:              "Brazier",
		Long:               `Brazier`,
		Run:                a.Run,
		SilenceErrors:      true,
		SilenceUsage:       true,
		PersistentPostRunE: a.Cleanup,
	}

	cmd.SetOutput(os.Stdout)
	cmd.AddCommand(NewCreateCmd(&a))

	return &cmd
}

// App is the main cli application
type app struct {
	Path string
	DB   *storm.DB
	Out  io.Writer
}

// Registrar returns the active registrar
func (a *app) db() (*storm.DB, error) {
	if a.DB == nil {
		var err error

		// Using BoltDB for now
		// In the future, use config
		a.DB, err = storm.Open(a.Path, storm.AutoIncrement())
		if err != nil {
			return nil, err
		}
	}

	return a.DB, nil
}

// Registrar returns the active registrar
func (a *app) Registrar() (brazier.Registrar, error) {
	db, err := a.db()
	if err != nil {
		return nil, err
	}

	return boltdb.NewRegistrar(db), nil
}

// Registrar returns the active registrar
func (a *app) Store() (brazier.Store, error) {
	db, err := a.db()
	if err != nil {
		return nil, err
	}

	return boltdb.NewStore(db), nil
}

// Run runs the root command
func (a *app) Run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// Cleanup is run after the command
func (a *app) Cleanup(cmd *cobra.Command, args []string) error {
	if a.DB != nil {
		err := a.DB.Close()
		if err != nil {
			return err
		}

		a.DB = nil
		return nil
	}

	return nil
}
