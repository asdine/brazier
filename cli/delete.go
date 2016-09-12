package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewDeleteCmd creates a "Delete" cli command
func NewDeleteCmd(a *app) *cobra.Command {
	deleteCmd := deleteCmd{App: a}

	cmd := cobra.Command{
		Use:   "delete",
		Short: "Delete a value from a bucket",
		Long:  `Delete a value from a bucket`,
		RunE:  deleteCmd.Delete,
	}

	return &cmd
}

type deleteCmd struct {
	App *app
}

func (d *deleteCmd) Delete(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Wrong number of arguments")
	}

	bucket, err := d.App.Store.Bucket(args[0])
	if err != nil {
		return err
	}
	defer bucket.Close()

	err = bucket.Delete(args[1])
	if err != nil {
		return err
	}

	fmt.Fprintf(d.App.Out, "Item \"%s\" successfully deleted.\n", args[1])
	return nil
}
