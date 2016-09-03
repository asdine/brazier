package cli

import (
	"fmt"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/rpc"
	"github.com/spf13/cobra"
)

// NewRPCCmd creates a "RPC" cli command
func NewRPCCmd(a *app) *cobra.Command {
	rpcCmd := rpcCmd{
		App:       a,
		ServeFunc: rpc.Serve,
	}

	cmd := cobra.Command{
		Use:   "rpc",
		Short: "Runs Brazier as an RPC server",
		Long:  "Runs Brazier as an RPC server",
		RunE:  rpcCmd.Serve,
	}

	cmd.Flags().IntVarP(&rpcCmd.App.Config.RPC.Port, "port", "p", 5657, "Port")
	return &cmd
}

type rpcCmd struct {
	App       *app
	Port      int
	ServeFunc func(brazier.Store, int) error
}

func (r *rpcCmd) Serve(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(r.App.Out, "Serving RPC on port %d\n", r.App.Config.RPC.Port)

	return r.ServeFunc(r.App.Store, r.App.Config.RPC.Port)
}
