package cli

import (
	"errors"
	"fmt"

	"github.com/asdine/brazier/json"
	"github.com/spf13/cobra"
)

// NewListCmd creates a "List" cli command
func NewListCmd(a *app) *cobra.Command {
	listCmd := listCmd{App: a}

	cmd := cobra.Command{
		Use:   "list",
		Short: "Lists the content of a bucket",
		Long:  `Lists the content of a bucket`,
		RunE:  listCmd.List,
	}

	return &cmd
}

type listCmd struct {
	App *app
}

func (l *listCmd) List(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong number of arguments")
	}

	bucket, err := l.App.Store.Bucket(args[0])
	if err != nil {
		return err
	}
	defer bucket.Close()

	items, err := bucket.Page(1, -1)
	if err != nil {
		return err
	}

	raw, err := json.MarshalList(items)
	if err != nil {
		return err
	}

	fmt.Fprintln(l.App.Out, string(raw))
	return nil
}
