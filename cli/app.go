package cli

import (
	"io"
	"os"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store/boltdb"
	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New() *cobra.Command {
	a := app{Out: os.Stdout, Path: "brazier.db"}

	cmd := cobra.Command{
		Use:           "brazier",
		Short:         "Brazier",
		Long:          `Brazier`,
		Run:           a.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.SetOutput(os.Stdout)
	cmd.AddCommand(NewCreateCmd(&a))
	cmd.AddCommand(NewSaveCmd(&a))
	cmd.AddCommand(NewGetCmd(&a))
	cmd.AddCommand(NewHTTPCmd(&a))

	return &cmd
}

// App is the main cli application
type app struct {
	Path string
	Out  io.Writer
}

// Store returns the boltdb store
func (a *app) Store() (brazier.Store, error) {
	return boltdb.NewStore(a.Path), nil
}

// Run runs the root command
func (a *app) Run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
