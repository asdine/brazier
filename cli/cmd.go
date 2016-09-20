package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/asdine/brazier/json"
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
	cmd.AddCommand(NewUseCmd(&a))

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
			if len(args) != 2 {
				return errors.New("Wrong number of arguments")
			}

			name, err := a.defaultBucket()
			if err != nil {
				return err
			}

			err = a.Cli.Save(name, args[0], []byte(args[1]))
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
func NewGetCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get a value from a bucket",
		Long:  `Get a value in a bucket`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Wrong number of arguments")
			}

			name, err := a.defaultBucket()
			if err != nil {
				return err
			}

			out, err := a.Cli.Get(name, args[0])
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
		Long:  "Lists all the buckets.\nIf a bucket name is specified, lists the content of the bucket instead.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("Wrong number of arguments")
			}

			if len(args) == 0 {
				list, err := a.Cli.ListBuckets()
				if err != nil {
					return err
				}

				deflt, err := a.defaultBucket()
				if err != nil {
					return err
				}

				for i := range list {
					a.Out.Write([]byte(list[i]))
					if list[i] == deflt {
						a.Out.Write([]byte(" *"))
					}
					a.Out.Write([]byte("\n"))
				}
				return nil
			}

			items, err := a.Cli.List(args[0])
			if err != nil {
				return err
			}

			data, err := json.MarshalList(items)
			if err != nil {
				return err
			}

			a.Out.Write(data)
			a.Out.Write([]byte("\n"))
			return nil
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
			if len(args) != 1 {
				return errors.New("Wrong number of arguments")
			}

			name, err := a.defaultBucket()
			if err != nil {
				return err
			}

			err = a.Cli.Delete(name, args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(a.Out, "Item \"%s\" successfully deleted.\n", args[0])
			return nil
		},
	}

	return &cmd
}

// NewUseCmd creates a "Use" cli command
func NewUseCmd(a *app) *cobra.Command {
	cmd := cobra.Command{
		Use:   "use",
		Short: "Set a bucket as default",
		Long:  `Set a bucket as default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Bucket name is missing")
			}

			if args[0] != defaultBucket {
				names, err := a.Cli.ListBuckets()
				if err != nil {
					return err
				}

				var found bool
				for _, name := range names {
					if name == args[0] {
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("Bucket \"%s\" not found.\n", args[0])
				}
			}

			db, err := a.settingsDB()
			if err != nil {
				return err
			}
			defer db.Close()

			return db.Set("buckets", "default", args[0])
		},
	}

	return &cmd
}
