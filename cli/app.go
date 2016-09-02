package cli

import (
	"io"
	"os"

	"github.com/asdine/brazier"
	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New(s brazier.Store) *cobra.Command {
	a := app{Out: os.Stdout, Store: s}

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
	cmd.AddCommand(NewDeleteCmd(&a))
	cmd.AddCommand(NewListCmd(&a))
	cmd.AddCommand(NewHTTPCmd(&a))
	cmd.AddCommand(NewRPCCmd(&a))

	return &cmd
}

// App is the main cli application
type app struct {
	Out   io.Writer
	Store brazier.Store
}

// Run runs the root command
func (a *app) Run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
