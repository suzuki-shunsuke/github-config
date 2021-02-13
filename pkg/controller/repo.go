package controller

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	maxPerPage = 100
	maxLoop    = 100
)

type RepoEditor interface {
	Edit(*github.Repository)
}

func (ctrl *Controller) RunRepo(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"count": len(cfg.Repo.Rules),
	}).Info("list repo rules")

	ctrl.Config = cfg
	param.Owner = cfg.Owner

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: param.GitHubToken},
	)))
	for i, rule := range cfg.Repo.Rules {
		rule.Policy.SetGitHubClient(client)
		cfg.Repo.Rules[i] = rule
	}

	repos, err := ctrl.listRepos(ctx, client, param)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"count": len(repos),
	}).Info("list repositories")
	for _, repo := range repos {
		r := Repository{
			GitHub: repo,
			Name:   repo.GetName(),
			Owner:  cfg.Owner,
		}
		logE := logrus.WithFields(logrus.Fields{
			"repo": r.Name,
		})
		if err := ctrl.handleRepo(ctx, param, client, r); err != nil {
			logE.WithError(err).Error("handle repository")
		}
	}
	return nil
}

func (ctrl *Controller) listRepos(ctx context.Context, client *github.Client, param Param) ([]*github.Repository, error) {
	var arr []*github.Repository
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: maxPerPage,
		},
	}
	for i := 0; i < maxLoop; i++ {
		repos, _, err := client.Repositories.List(ctx, param.Owner, opt)
		if err != nil {
			return nil, fmt.Errorf("list repositories (owner: %s, page: %d): %w", param.Owner, opt.Page, err)
		}
		arr = append(arr, repos...)
		if len(repos) != maxPerPage {
			return arr, nil
		}
		opt.Page += 1
	}
	return arr, nil
}

type Repository struct {
	GitHub *github.Repository
	Owner  string
	Name   string
}

func (ctrl *Controller) handleRepo(ctx context.Context, param Param, client *github.Client, repo Repository) error { //nolint:unparam
	repoName := repo.Name
	logE := logrus.WithFields(logrus.Fields{
		"repo": repoName,
	})
	updatedRepo := &github.Repository{}

	isEdited := false
	for _, rule := range ctrl.Config.Repo.Rules {
		logE.WithFields(logrus.Fields{
			"targets": rule.Targets,
		}).Debug("check rule")
		if f, err := rule.Targets.Match(repo); err != nil { //nolint:nestif
			logE.WithError(err).Error("check a repository matches with the targets")
			continue
		} else if f {
			logE.Debug("a repository matches with the targets")
			if f, err := rule.Policy.Match(ctx, repo); err != nil {
				logE.WithError(err).Error("check a repository matches with the policy")
				continue
			} else if f {
				logE.Info("a repository matches with the rule")
				if edit, ok := rule.Policy.(RepoEditor); ok {
					isEdited = true
					logE.Info("edit a repository")
					edit.Edit(updatedRepo)
				}
			}
		}
	}
	if isEdited {
		if param.DryRun {
			logE.Info("[DRY RUN] update a repository")
		} else {
			if _, _, err := client.Repositories.Edit(ctx, param.Owner, repoName, updatedRepo); err != nil {
				logE.WithError(err).Error("update a repository")
			}
			logE.Info("update a repository")
		}
	}
	return nil
}
