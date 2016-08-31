package cli

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/http"
	"github.com/spf13/cobra"
)

// NewHTTPCmd creates a "HTTP" cli command
func NewHTTPCmd(a *app) *cobra.Command {
	httpCmd := httpCmd{
		App:       a,
		ServeFunc: http.Serve,
	}

	cmd := cobra.Command{
		Use:   "http",
		Short: "Runs Brazier as an HTTP server",
		Long:  "Runs Brazier as an HTTP server",
		RunE:  httpCmd.Serve,
	}

	cmd.Flags().IntVarP(&httpCmd.Port, "port", "p", 5656, "Port")

	return &cmd
}

type httpCmd struct {
	App       *app
	Port      int
	ServeFunc func(brazier.Store, int) error
}

func (h *httpCmd) Serve(cmd *cobra.Command, args []string) error {
	err := h.ServeFunc(h.App.Store, h.Port)
	if err != nil {
		return err
	}
	return nil
}
