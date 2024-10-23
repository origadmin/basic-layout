// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package start is the start command for the application.
package start

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	logger "github.com/origadmin/slog-kratos"
	"github.com/spf13/cobra"

	"origadmin/basic-layout/internal/bootstrap"
	"origadmin/basic-layout/internal/mods"
	"origadmin/basic-layout/toolkits/utils"
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
	name = "origadmin.server.v1"
	// Version is the version of the compiled software.
	version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.DefaultFlags()
)

var cmd = &cobra.Command{
	Use:   "start",
	Short: "start the server",
	RunE:  startRun,
}

func init() {
	flags.Name = name
	flags.Version = version

}

// Cmd The function defines a CLI command to start a server with various flags and options, including the
// ability to run as a daemon.
func Cmd() *cobra.Command {
	cmd.Flags().BoolP(startRandom, "r", false, "start with random password")
	cmd.Flags().StringP(startWorkDir, "d", ".", "working directory")
	cmd.Flags().StringP(startConfig, "c", "resources",
		"runtime configuration files or directory (relative to workdir, multiple separated by commas)")
	cmd.Flags().StringP(startStatic, "s", "", "static files directory")
	cmd.Flags().Bool(startDaemon, false, "run as a daemon")
	return cmd
}

func startRun(cmd *cobra.Command, args []string) error {
	flags.WorkDir, _ = cmd.Flags().GetString(startWorkDir)
	staticDir, _ := cmd.Flags().GetString(startStatic)
	flags.ConfigPath, _ = cmd.Flags().GetString(startConfig)
	//random, _ := cmd.Flags().GetBool(startRandom)

	flags.MetaData = make(map[string]string)
	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.ID,
		"service.name", flags.Name,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(l)
	path := filepath.Join(flags.WorkDir, flags.ConfigPath)
	//envpath := filepath.Join(flags.WorkDir, flags.EnvPath)
	log.Infow(startWorkDir, flags.WorkDir, startStatic, staticDir, startConfig, path)
	//env, _ := bootstrap.LoadEnv(envpath)
	bs, err := bootstrap.FromLocal(name, path, nil, l)
	if err != nil {
		return err
	}

	if daemon, _ := cmd.Flags().GetBool("daemon"); daemon {
		bin, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Errorf("failed to get absolute path for command: %s \n", err.Error())
			return err
		}

		cmdArgs := []string{"start"}
		cmdArgs = append(cmdArgs, "-d", strings.TrimSpace(flags.WorkDir))
		cmdArgs = append(cmdArgs, "-c", strings.TrimSpace(flags.ConfigPath))
		cmdArgs = append(cmdArgs, "-s", strings.TrimSpace(staticDir))
		_, _ = fmt.Printf("execute command: %s %s \n", bin, strings.Join(cmdArgs, " "))
		command := exec.Command(bin, cmdArgs...)
		err = command.Start()
		if err != nil {
			_, _ = fmt.Printf("failed to start daemon thread: %s \n", err.Error())
			return err
		}

		pid := command.Process.Pid
		//err = os.WriteFile(
		//	fmt.Sprintf("%s.lock", utils.ToLower(cmd)),
		//	[]byte(fmt.Sprintf("%d", pid)),
		//	0o600)
		//if err != nil {
		//	log.Errorf("failed to write pid file: %s \n", err.Error())
		//}
		log.Errorf("service %s daemon thread started with pid %d \n", bs.ServiceName, pid)
		return nil
	}
	lockfile := fmt.Sprintf("%s.lock", utils.ToLower(cmd))
	err = os.WriteFile(
		lockfile,
		[]byte(fmt.Sprintf("%d", os.Getpid())),
		0o600)
	if err == nil {
		defer os.Remove(lockfile)
	}
	//info to ctx
	app, cleanup, err := buildInjectors(cmd.Context(), bs, l)
	if err != nil {
		return err
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func NewApp(ctx context.Context, injector *mods.InjectorClient) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(flags.ID),
		kratos.Name(flags.Name),
		kratos.Version(flags.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		//kratos.Server(hs, gs, gss),
		//kratos.Server(injector.ServerGINS),
	}

	bootstrap.InjectorGinServer(injector)
	if injector.ServerGINS != nil {
		opts = append(opts, kratos.Server(injector.ServerGINS))
	}

	return kratos.New(opts...)
}
