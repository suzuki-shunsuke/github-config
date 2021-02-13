package controller

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
)

const DataDogEventSourceType = "GITHUBCONFIG"

type RuleHasProjects struct {
	client            *github.Client
	CheckListProjects bool
}

func (rule *RuleHasProjects) SetGitHubClient(client *github.Client) {
	rule.client = client
}

func (rule *RuleHasProjects) Edit(updatedRepo *github.Repository) {
	updatedRepo.HasProjects = github.Bool(false)
}

func (rule *RuleHasProjects) Match(ctx context.Context, repo Repository) (bool, error) {
	if !repo.GitHub.GetHasProjects() {
		return false, nil
	}
	logE := logrus.WithFields(logrus.Fields{
		"repo": repo.Name,
	})
	if !rule.CheckListProjects {
		return true, nil
	}
	if projects, _, err := rule.client.Repositories.ListProjects(ctx, repo.Owner, repo.Name, nil); err != nil {
		logE.WithError(err).Error("list projects")
		return true, nil
	} else if len(projects) == 0 {
		return true, nil
	}
	return false, nil
}
