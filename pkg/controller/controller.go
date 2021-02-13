package controller

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

type Controller struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Config  Config
	DataDog *datadog.Client
	Rules   []Rule
}

func New(ctx context.Context, param Param) (Controller, Param, error) {
	if param.LogLevel != "" {
		lvl, err := logrus.ParseLevel(param.LogLevel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"log_level": param.LogLevel,
			}).WithError(err).Error("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	return Controller{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, param, nil
}
