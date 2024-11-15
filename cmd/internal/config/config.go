// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package start is the start command for the application.
package start

import (
	"github.com/spf13/cobra"

	"origadmin/basic-layout/internal/bootstrap"
)

const (
	startRandom  = `random`
	startWorkDir = `workdir`
	startConfig  = `config`
	startStatic  = `static`
	startDaemon  = `daemon`
)

var (
	// Name is the name of the compiled software.
	Name = "origadmin.server.v1"
	// Version is the Version of the compiled software.
	Version = "v1.0.0"
	// boot are the bootstrap boot.
	boot = bootstrap.Bootstrap{}
)

var cmd = &cobra.Command{
	Use:   "config",
	Short: "the config command for the application",
	RunE:  configCommandRun,
}

func init() {
	boot.SetFlags(Name, Version)
}

// Cmd The function defines a CLI command to start a server with various boot and options, including the
// ability to run as a daemon.
func Cmd() *cobra.Command {
	cmd.Flags().BoolP(startRandom, "r", false, "start with random password")
	cmd.Flags().StringP(startWorkDir, "d", ".", "working directory")
	cmd.Flags().StringP(startConfig, "c", "bootstrap.toml",
		"runtime configuration files or directory (relative to workdir, multiple separated by commas)")
	cmd.Flags().StringP(startStatic, "s", "", "static files directory")
	cmd.Flags().Bool(startDaemon, false, "run as a daemon")
	return cmd
}

func configCommandRun(cmd *cobra.Command, args []string) error {
	cmd.Println("config command run")
	return nil
}
