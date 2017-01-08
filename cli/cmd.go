package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// New returns a configured Cobra command
func New() *cobra.Command {
	a := app{
		Out: os.Stdout,
	}

	cmd := cobra.Command{
		Use:                "brazier",
		Short:              "A JSON storage with command line, HTTP and gRPC support.",
		Long:               `A JSON storage with command line, HTTP and gRPC support.`,
		Run:                a.Run,
		SilenceErrors:      true,
		SilenceUsage:       true,
		PersistentPreRunE:  a.PreRun,
		PersistentPostRunE: a.PostRun,
	}

	cmd.SetOutput(os.Stdout)
	cmd.AddCommand(NewCreateCmd(&a))
	cmd.AddCommand(NewPutCmd(&a))
	cmd.AddCommand(NewGetCmd(&a, false))
	cmd.AddCommand(NewDeleteCmd(&a))
	cmd.AddCommand(NewServerCmd(&a))

	cmd.PersistentFlags().StringVar(&a.ConfigPath, "config", "", "config file")
	cmd.PersistentFlags().StringVar(&a.DataDir, "data-dir", "", "data directory (default $HOME/.brazier)")
	return &cmd
}

// NewCreateCmd creates a "create" cli command
func NewCreateCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "create PATH",
		Short: "Create a bucket",
		Long: `Create a bucket at the given path.
A path is a bucket name or a list of bucket names separated by the character '/'.`,
		Example: `brazier create friends
brazier create food/vegetables
brazier create food/drinks/sodas`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("Bucket name is missing")
			}

			err := a.Cli.Create(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Bucket \"%s\" successfully created.\n", args[0])
			return nil
		},
	}

	return &cmd
}

// NewSaveCmd creates a "Save" cli command
func NewPutCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "put PATH",
		Short: "Set or replace a value in a bucket",
		Long: `Set or replace a value in a bucket. A value can be anything.
JSON values are automatically detected.`,
		Example: `brazier put friends/john/phone 555-666
brazier put users/1 '{"username": "john"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Wrong number of arguments")
			}

			err := a.Cli.Put(args[0], []byte(args[1]))
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Item \"%s\" successfully saved.\n", args[0])
			return nil
		},
	}

	return &cmd
}

// NewGetCmd creates a "Get" cli command
func NewGetCmd(a *app, recByDefault bool) *cobra.Command {
	var recursive bool

	cmd := cobra.Command{
		Use:   "get PATH",
		Short: "Get a value from a key or list bucket content",
		Long:  `Get a value from a key or list bucket content.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Wrong number of arguments")
			}

			out, err := a.Cli.Get(args[0], recursive)
			if err != nil {
				return err
			}

			_, err = a.Out.Write(out)
			return err
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", recByDefault, "display all the items recursively from the given bucket path.")

	return &cmd
}

// NewDeleteCmd creates a "Delete" cli command
func NewDeleteCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete PATH",
		Short: "Delete a key from a bucket",
		Long:  `Delete a key from a bucket.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Wrong number of arguments")
			}

			err := a.Cli.Delete(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Item \"%s\" successfully deleted.\n", args[0])
			return nil
		},
	}

	return &cmd
}
