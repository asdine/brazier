package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewGetCmd creates a "Get" cli command
func NewGetCmd(a *app) *cobra.Command {
	getCmd := getCmd{App: a}

	cmd := cobra.Command{
		Use:   "get",
		Short: "Get a value from a bucket",
		Long:  `Get a value in a bucket`,
		RunE:  getCmd.Get,
	}

	return &cmd
}

type getCmd struct {
	App *app
}

func (g *getCmd) Get(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Wrong number of arguments")
	}

	bucket, err := g.App.Store.Bucket(args[0])
	if err != nil {
		return err
	}
	defer bucket.Close()

	item, err := bucket.Get(args[1])
	if err != nil {
		return err
	}

	fmt.Fprintln(g.App.Out, string(item.Data))
	return nil
}
