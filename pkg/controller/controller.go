package controller

import (
	"context"
	"errors"
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
			return Controller{}, param, errors.New("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	param.GitHubToken = os.Getenv("GITHUB_TOKEN")
	if param.GitHubToken == "" {
		return Controller{}, param, errors.New("GITHUB_TOKEN is missing")
	}

	param.DataDogAPIKey = os.Getenv("DATADOG_API_KEY")

	return Controller{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, param, nil
}
