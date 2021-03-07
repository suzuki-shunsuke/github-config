package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"golang.org/x/oauth2"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const (
	maxPerPage = 100
	maxLoop    = 100
)

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

	if param.DataDogAPIKey != "" {
		ctrl.DataDog = datadog.NewClient(param.DataDogAPIKey, "")
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: param.GitHubToken},
	)))
	for i, rule := range cfg.Repo.Rules {
		rule.Policy.SetGitHubClient(client)
		rule.Policy.SetDataDogClient(ctrl.DataDog)
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
		r := domain.Repository{
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
		opt.Page++
	}
	return arr, nil
}

func (ctrl *Controller) repoAction(ctx context.Context, param *domain.ParamAction, policy domain.RepoPolicy) error {
	switch t := policy.Action().Type; t {
	case "datadog_metric":
	case "fix":
		a, ok := policy.(domain.Fixable)
		if !ok {
			return errors.New("this rule doesn't support to fix")
		}
		a.Fix(ctx, param)
	default:
		return errors.New("invalid action type: " + t)
	}
	return nil
}

func (ctrl *Controller) handleRepo(ctx context.Context, param Param, client *github.Client, repo domain.Repository) error { //nolint:unparam,cyclop
	repoName := repo.Name
	logE := logrus.WithFields(logrus.Fields{
		"repo": repoName,
	})
	ts := time.Now().Unix()
	paramAction := domain.ParamAction{
		Repo:             repo,
		UpdatedRepo:      &github.Repository{},
		TimestampFloat64: float64(ts),
		TimestampInt:     int(ts),
		DryRun:           param.DryRun,
	}
	for _, rule := range ctrl.Config.Repo.Rules {
		actionConfig := rule.Policy.Action()
		logE.WithFields(logrus.Fields{
			"target": rule.Target,
		}).Debug("check rule")
		if f, err := rule.Target.Match(repo); err != nil { //nolint:nestif
			logE.WithError(err).Error("check a repository matches with the targets")
			continue
		} else if f {
			logE.Debug("a repository matches with the targets")
			if actionConfig.Type == "datadog_metric" {
				paramAction.DataDogMetrics = append(paramAction.DataDogMetrics, rule.Policy.DataDogMetric(paramAction.Repo, &paramAction.TimestampFloat64))
			}
			if f, err := rule.Policy.Match(ctx, repo); err != nil {
				logE.WithError(err).Error("check a repository matches with the policy")
				continue
			} else if f {
				logE.Info("a repository matches with the rule")
				// TODO action
				if err := ctrl.repoAction(ctx, &paramAction, rule.Policy); err != nil {
					logE.WithError(err).Error("prepare")
					continue
				}
			}
		}
	}
	if paramAction.IsEdited {
		if param.DryRun {
			logE.Info("[DRY RUN] update a repository")
		} else {
			if _, _, err := client.Repositories.Edit(ctx, param.Owner, repoName, paramAction.UpdatedRepo); err != nil {
				logE.WithError(err).Error("update a repository")
			}
			logE.Info("update a repository")
		}
	}
	if len(paramAction.DataDogMetrics) != 0 {
		if param.DryRun {
			logE.Info("[DRY RUN] post metrics to DataDog")
		} else {
			if err := ctrl.DataDog.PostMetrics(paramAction.DataDogMetrics); err != nil {
				logE.WithError(err).Error("post metrics to DataDog")
			}
			logE.Info("post metrics to DataDog")
		}
	}
	return nil
}
