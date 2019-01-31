package container

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type RestartOptions struct {
	NSeconds        int
	NSecondsChanged bool

	Containers []string
}

// NewRestartCommand creates a new cobra.Command for `docker restart`
func NewRestartCommand(dockerCli command.Cli) *cobra.Command {
	var opts RestartOptions

	cmd := &cobra.Command{
		Use:   "restart [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Restart one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Containers = args
			opts.NSecondsChanged = cmd.Flags().Changed("time")
			return RunRestart(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.NSeconds, "time", "t", 10, "Seconds to wait for stop before killing the container")
	return cmd
}

func RunRestart(dockerCli command.Cli, opts *RestartOptions) error {
	ctx := context.Background()
	var errs []string
	var timeout *time.Duration
	if opts.NSecondsChanged {
		timeoutValue := time.Duration(opts.NSeconds) * time.Second
		timeout = &timeoutValue
	}

	for _, name := range opts.Containers {
		if err := dockerCli.Client().ContainerRestart(ctx, name, timeout); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintln(dockerCli.Out(), name)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
