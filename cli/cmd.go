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
		Short:              "Brazier",
		Long:               `Brazier`,
		Run:                a.Run,
		SilenceErrors:      true,
		SilenceUsage:       true,
		PersistentPreRunE:  a.PreRun,
		PersistentPostRunE: a.PostRun,
	}

	cmd.SetOutput(os.Stdout)
	cmd.AddCommand(NewCreateCmd(&a))
	cmd.AddCommand(NewSaveCmd(&a))
	cmd.AddCommand(NewGetCmd(&a))
	cmd.AddCommand(NewDeleteCmd(&a))
	cmd.AddCommand(NewListCmd(&a))
	cmd.AddCommand(NewServerCmd(&a))

	cmd.PersistentFlags().StringVar(&a.ConfigPath, "config", "", "config file")
	cmd.PersistentFlags().StringVar(&a.DataDir, "data-dir", "", "data directory (default $HOME/.brazier)")
	return &cmd
}

// NewCreateCmd creates a "create" cli command
func NewCreateCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "create",
		Short: "Create a bucket",
		Long:  `Create a bucket`,
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
func NewSaveCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "save",
		Short: "Save a value in a bucket",
		Long:  `Save a value in a bucket`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("Wrong number of arguments")
			}

			err := a.Cli.Save(args[0], args[1], []byte(args[2]))
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Item \"%s\" successfully saved.\n", args[1])
			return nil
		},
	}

	return &cmd
}

// NewGetCmd creates a "Get" cli command
func NewGetCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get a value from a bucket",
		Long:  `Get a value in a bucket`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Wrong number of arguments")
			}

			out, err := a.Cli.Get(args[0], args[1])
			if err != nil {
				return err
			}

			_, err = a.Out.Write(out)
			return err
		},
	}

	return &cmd
}

// NewListCmd creates a "List" cli command
func NewListCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List buckets or bucket content",
		Long:  "Lists all the buckets.\nIf a bucket name is specified, list the bucket's content instead.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("Wrong number of arguments")
			}

			var out []byte
			var err error

			if len(args) == 0 {
				out, err = a.Cli.ListBuckets()
			} else {
				out, err = a.Cli.List(args[0])
			}

			if err != nil {
				return err
			}

			_, err = a.Out.Write(out)
			return err
		},
	}

	return &cmd
}

// NewDeleteCmd creates a "Delete" cli command
func NewDeleteCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete",
		Short: "Delete a key from a bucket",
		Long:  `Delete a key from a bucket`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Wrong number of arguments")
			}

			err := a.Cli.Delete(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Item \"%s\" successfully deleted.\n", args[1])
			return nil
		},
	}

	return &cmd
}
