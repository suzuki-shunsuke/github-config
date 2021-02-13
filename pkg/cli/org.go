package cli

import (
	"fmt"

	"github.com/suzuki-shunsuke/github-config/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) orgAction(c *cli.Context) error {
	param, err := runner.setCLIArg(c, controller.Param{})
	if err != nil {
		return fmt.Errorf("parse the command line arguments: %w", err)
	}

	ctrl, param, err := controller.New(c.Context, param)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.RunOrg(c.Context, param)
}
