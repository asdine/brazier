package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewCreateCmd creates a "create" cli command
func NewCreateCmd(a *app) *cobra.Command {
	createCmd := createCmd{App: a}

	cmd := cobra.Command{
		Use:   "create",
		Short: "Create a bucket",
		Long:  `Create a bucket`,
		RunE:  createCmd.Create,
	}

	return &cmd
}

type createCmd struct {
	App *app
}

func (c *createCmd) Create(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Bucket name is missing")
	}

	err := c.App.Store.Create(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(c.App.Out, "Bucket \"%s\" successfully created.\n", args[0])
	return nil
}
