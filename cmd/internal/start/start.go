// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package start is the start command for the application.
package start

import (
	"github.com/spf13/cobra"
)

const (
	startRandom  = `random`
	startWorkDir = `workdir`
	startConfig  = `config`
	startStatic  = `static`
	startDaemon  = `daemon`
)

// Cmd The function defines a CLI command to start a server with various flags and options, including the
// ability to run as a daemon.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
