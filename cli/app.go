package cli

import (
	"io"
	"os"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/config"
	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New(s brazier.Store) *cobra.Command {
	a := app{Out: os.Stdout, Store: s}

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
	return &cmd
}

// App is the main cli application
type app struct {
	Out        io.Writer
	Store      brazier.Store
	ConfigPath string
	Config     config.Config
}

// Run runs the root command
func (a *app) Run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// PreRun runs the root command
func (a *app) PreRun(cmd *cobra.Command, args []string) error {
	if a.ConfigPath != "" {
		return config.FromFile(a.ConfigPath, &a.Config)
	}
	return nil
}
