package controller

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	maxPerPage = 100
	maxLoop    = 100
)

func (ctrl *Controller) Run(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)))

	org, _, err := client.Organizations.Get(ctx, param.Owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"owner": param.Owner,
		}).WithError(err).Error("get an organization")
	} else if err := ctrl.handleOrg(ctx, param, client, org); err != nil {
		logrus.WithError(err).Error("update org")
	}

	repos, err := ctrl.listRepos(ctx, client, param)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		repoName := repo.GetName()
		logE := logrus.WithFields(logrus.Fields{
			"repo": repoName,
		})
		if err := ctrl.handleRepo(ctx, param, client, repo); err != nil {
			logE.WithError(err).Error("handle repository")
		}
	}
	return nil
}
