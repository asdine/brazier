package cli

import (
	"fmt"

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
		Short: "Run Brazier as an HTTP server",
		Long:  "Run Brazier as an HTTP server",
		RunE:  httpCmd.Serve,
	}

	cmd.Flags().IntVarP(&httpCmd.App.Config.HTTP.Port, "port", "p", 5656, "Port")
	return &cmd
}

type httpCmd struct {
	App       *app
	ServeFunc func(brazier.Store, int) error
}

func (h *httpCmd) Serve(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(h.App.Out, "Serving HTTP on port %d\n", h.App.Config.HTTP.Port)

	return h.ServeFunc(h.App.Store, h.App.Config.HTTP.Port)
}
