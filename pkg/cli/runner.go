package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/github-config/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "github-config",
		Usage:   "GitHub Organization Configuration Management. https://github.com/suzuki-shunsuke/github-config",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "repo",
				Usage:  "GitHub Repositories Configuration Management",
				Action: runner.repoAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "configuration file path",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "dry run",
					},
				},
			},
			{
				Name:   "org",
				Usage:  "GitHub Organization Configuration Management",
				Action: runner.orgAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "configuration file path",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "dry run",
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}
