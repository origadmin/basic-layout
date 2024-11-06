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
	"github.com/origadmin/toolkits/errors"
	"github.com/spf13/cobra"

	"origadmin/basic-layout/internal/bootstrap"
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
	Name = "origadmin.server.v1"
	// Version is the Version of the compiled software.
	Version = "v1.0.0"
	// flags are the bootstrap flags.
	flags = bootstrap.BootFlags{}
)

var cmd = &cobra.Command{
	Use:   "start",
	Short: "start the server",
	RunE:  startRun,
}

func init() {
	flags = bootstrap.NewBootFlags(Name, Version)
}

// Cmd The function defines a CLI command to start a server with various flags and options, including the
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

func startRun(cmd *cobra.Command, args []string) error {
	flags.WorkDir, _ = cmd.Flags().GetString(startWorkDir)
	staticDir, _ := cmd.Flags().GetString(startStatic)
	flags.ConfigPath, _ = cmd.Flags().GetString(startConfig)
	//random, _ := cmd.Flags().GetBool(startRandom)

	flags.Metadata = make(map[string]string)
	l := log.With(logger.NewLogger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", flags.ID,
		"service.name", flags.ServiceName,
		"service.version", flags.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(l)
	//path := filepath.Join(flags.WorkDir, flags.ConfigPath)
	//envpath := filepath.Join(flags.WorkDir, flags.EnvPath)
	log.Infow("msg", "start info", startWorkDir, flags.WorkDir, startStatic, staticDir, startConfig, flags.ConfigPath)
	//env, _ := bootstrap.LoadEnv(envpath)
	//bs, err := bootstrap.FromLocalPath(flags.ServiceName, path, l)
	//if err != nil {
	//	return errors.Wrap(err, "load config error")
	//}
	src := bootstrap.LoadSourceFiles(flags.WorkDir, flags.ConfigPath)
	bs, err := bootstrap.FromRemote(flags.ServiceName, src)
	if err != nil {
		return errors.Wrap(err, "load config error")
	}
	if bs == nil {
		return fmt.Errorf("bootstrap config not found")
	}

	log.Infof("bootstrap: %+v", bootstrap.PrintString(bs))

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
		log.Errorf("service %s daemon thread started with pid %d \n", bs.ServiceName, pid)
		return nil
	}
	lockfile := fmt.Sprintf("%s.lock", utils.ToLower(cmd))
	if err = os.WriteFile(lockfile, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600); err == nil {
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

func NewApp(ctx context.Context, injector *bootstrap.InjectorClient) *kratos.App {
	opts := []kratos.Option{
		kratos.ID(flags.ID),
		kratos.Name(flags.ServiceName),
		kratos.Version(flags.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Context(ctx),
		kratos.Signal(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT),
		kratos.Logger(injector.Logger),
		//kratos.Server(hs, gs, gss),
		//kratos.Server(injector.ServerGINS),
	}

	err := bootstrap.InjectorGinServer(injector)
	if err != nil {
		panic(err)
	}
	//if injector.ServerGINS != nil {
	//	opts = append(opts, kratos.Server(injector.ServerGINS))
	//}
	if injector.ServerHTTP != nil {
		opts = append(opts, kratos.Server(injector.ServerHTTP))
	}

	return kratos.New(opts...)
}
