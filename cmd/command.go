package cmd

import (
	"github.com/spf13/cobra"

	"origadmin/basic-layout/cmd/internal/start"
)

func Commands() []*cobra.Command {
	return []*cobra.Command{
		start.StartCmd(),
	}
}
