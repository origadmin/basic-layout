package cmd

import (
	"github.com/spf13/cobra"

	"origadmin/basic-layout/cmd/server"
)

func Commands() []*cobra.Command {
	return []*cobra.Command{
		server.StartCmd(),
	}
}
