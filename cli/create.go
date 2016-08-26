package cli

import (
	"errors"
	"fmt"

	"github.com/asdine/brazier/store"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates a "create" cli command
func NewCreateCmd(a *app) *cobra.Command {
	createCmd := createCmd{App: a}

	cmd := cobra.Command{
		Use:   "create",
		Short: "Create a bucket",
		Long:  `Creates a bucket in Brazier`,
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

	reg, err := c.App.Registrar()
	if err != nil {
		return err
	}

	s, err := c.App.Store()
	if err != nil {
		return err
	}

	_, err = reg.Create(args[0], s)
	if err != nil {
		if err == store.ErrAlreadyExists {
			return fmt.Errorf("The bucket \"%s\" already exists.\n", args[0])
		}
		return err
	}

	fmt.Fprintf(c.App.Out, "Bucket \"%s\" successfully created.\n", args[0])
	return nil
}
