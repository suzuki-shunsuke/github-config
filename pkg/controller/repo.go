package controller

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (ctrl *Controller) RunRepo(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}

	ctrl.Config = cfg

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)))

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

type Rule interface {
	Match(repo Repository) (bool, error)
}

type ParamRule struct {
}

type EditRepo interface {
}

func (ctrl *Controller) handleRepo(ctx context.Context, param Param, client *github.Client, repo *github.Repository) error { //nolint:unparam
	repoName := repo.GetName()
	logE := logrus.WithFields(logrus.Fields{
		"repo": repoName,
	})
	updatedRepo := &github.Repository{}

	for _, rule := range ctrl.Rules {
		if f, err := rule.Match(ctx, repo); err != nil {
			logE.WithError(err).Error("update a repository")
			continue
		} else if f {
			if err := rule.Do(ctx, repo, updatedRepo); err != nil {
				logE.WithError(err).Error("update a repository")
			}
		}
	}
	if _, _, err := client.Repositories.Edit(ctx, param.Owner, repoName, updatedRepo); err != nil {
		logE.WithError(err).Error("update a repository")
	}

	for _, repoItem := range ctrl.Config.Repo.Items {
		if f, err := repoItem.Condition.MatchRepo(repo); err != nil {
			return err
		} else if !f {
			continue
		}
		if f, err := repoItem.Rule.MatchRepo(repo); err != nil {
			return err
		} else if !f {
			continue
		}
		if err := repoItem.DoAction(repo, updatedRepo); err != nil {
			return err
		}
	}

	if !repo.GetPrivate() {
		// public
		updatedRepo.Private = github.Bool(true)
	}
	if repo.GetHasProjects() {
		if projects, _, err := client.Repositories.ListProjects(ctx, param.Owner, repoName, nil); err != nil {
			logE.WithError(err).Error("list projects")
		} else if len(projects) == 0 {
			updatedRepo.HasProjects = github.Bool(false)
		}
	}
	if repo.GetHasIssues() {
		if issues, _, err := client.Issues.ListByRepo(ctx, param.Owner, repoName, nil); err != nil {
			logE.WithError(err).Error("list issues")
		} else if len(issues) == 0 {
			updatedRepo.HasIssues = github.Bool(false)
		}
	}
	if repo.GetHasPages() {
		// DisablePages
		updatedRepo.HasPages = github.Bool(false)
	}
	if repo.GetFork() {
		updatedRepo.Fork = github.Bool(false)
	}
	if _, _, err := client.Repositories.Edit(ctx, param.Owner, repoName, updatedRepo); err != nil {
		logE.WithError(err).Error("update a repository")
	}
	// UpdateBranchProtection
	if _, _, err := client.Repositories.UpdateBranchProtection(ctx, param.Owner, repoName, "master", &github.ProtectionRequest{}); err != nil {
		logE.WithError(err).Error("update a repository")
	}
	// TODO manage access
	// TODO remove not allowed webhook url
	// TODO remove not allowed integration
	// TODO remove not allowed deploy key
	if f, _, err := client.Repositories.GetVulnerabilityAlerts(ctx, param.Owner, repoName); err != nil {
		logE.WithError(err).Error("get whether the vulnerability alerts is enabled")
	} else if !f {
		if _, err := client.Repositories.EnableVulnerabilityAlerts(ctx, param.Owner, repoName); err != nil {
			logE.WithError(err).Error("enable vulnerability alerts")
		}
	}
	// TODO enable Dependabot security updates
	if _, err := client.Repositories.EnableAutomatedSecurityFixes(ctx, param.Owner, repoName); err != nil {
		logE.WithError(err).Error("enable automated security fixes")
	}

	return nil
}
