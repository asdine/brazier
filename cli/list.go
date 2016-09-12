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
		Short: "List buckets or bucket content",
		Long:  "Lists all the buckets.\nIf a bucket name is specified, list the bucket's content instead.",
		RunE:  listCmd.List,
	}

	return &cmd
}

type listCmd struct {
	App *app
}

func (l *listCmd) List(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Wrong number of arguments")
	}

	if len(args) == 0 {
		return l.listBuckets()
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

func (l *listCmd) listBuckets() error {
	list, err := l.App.Store.List()
	if err != nil {
		return err
	}

	for i := range list {
		fmt.Fprintf(l.App.Out, "%s\n", list[i])
	}

	return nil
}
