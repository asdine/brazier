package cli

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/asdine/brazier/json"
	"github.com/spf13/cobra"
)

// NewSaveCmd creates a "Save" cli command
func NewSaveCmd(a *app) *cobra.Command {
	saveCmd := saveCmd{App: a}

	cmd := cobra.Command{
		Use:   "save",
		Short: "Saves a value in a bucket",
		Long:  `Saves a value in a bucket`,
		RunE:  saveCmd.Save,
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

	bucket, err := s.App.Store.Bucket(args[0])
	if err != nil {
		return err
	}
	defer bucket.Close()

	data := []byte(args[2])
	if !json.IsValid(data) {
		var buffer bytes.Buffer
		buffer.Grow(len(data) + 2)
		buffer.WriteByte('"')
		buffer.Write(data)
		buffer.WriteByte('"')
		data = buffer.Bytes()
	} else {
		data = json.Clean(data)
	}

	_, err = bucket.Save(args[1], data)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.App.Out, "Item \"%s\" successfully saved.\n", args[1])
	return nil
}
