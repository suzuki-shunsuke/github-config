package cli

import (
	"fmt"

	"github.com/suzuki-shunsuke/github-config/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) setCLIArg(c *cli.Context, param controller.Param) (controller.Param, error) { //nolint:unparam
	if logLevel := c.String("log-level"); logLevel != "" {
		param.LogLevel = logLevel
	}
	param.ConfigFilePath = c.String("config")
	if param.ConfigFilePath == "" {
		param.ConfigFilePath = "github-config.yaml"
	}
	param.DryRun = c.Bool("dry-run")
	return param, nil
}

func (runner *Runner) repoAction(c *cli.Context) error {
	param, err := runner.setCLIArg(c, controller.Param{})
	if err != nil {
		return fmt.Errorf("parse the command line arguments: %w", err)
	}

	ctrl, param, err := controller.New(c.Context, param)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.RunRepo(c.Context, param)
}
