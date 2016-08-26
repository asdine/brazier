package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// NewSaveCmd creates a "Save" cli command
func NewSaveCmd(a *app) *cobra.Command {
	SaveCmd := saveCmd{App: a}

	cmd := cobra.Command{
		Use:   "save",
		Short: "Saves a value in a bucket",
		Long:  `Saves a value in a bucket`,
		RunE:  SaveCmd.Save,
	}

	return &cmd
}

type saveCmd struct {
	App *app
}

func (s *saveCmd) Save(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return errors.New("Wrong number of arguments")
	}

	store, err := s.App.Store()
	if err != nil {
		return err
	}

	bucket, err := store.Bucket(args[0])
	if err != nil {
		return err
	}
	defer bucket.Close()

	raw := []byte(args[2])

	var value interface{}
	err = json.Unmarshal(raw, &value)
	if err != nil {
		raw, err = json.Marshal(args[2])
		if err != nil {
			return err
		}
	}

	_, err = bucket.Save(args[1], raw)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.App.Out, "Item \"%s\" successfully saved.\n", args[1])
	return nil
}
