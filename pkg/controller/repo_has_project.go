package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v33/github"
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
	if !rule.CheckListProjects {
		return true, nil
	}
	if projects, _, err := rule.client.Repositories.ListProjects(ctx, repo.Owner, repo.Name, nil); err != nil {
		return false, fmt.Errorf("list projects: %w", err)
	} else if len(projects) == 0 {
		return true, nil
	}
	return false, nil
}

func newRuleHasProjects(param map[string]interface{}) (RepoPolicy, error) {
	policy := RuleHasProjects{}
	if a, ok := param["check_usage"]; !ok {
		return &RuleHasProjects{
			CheckListProjects: true,
		}, nil
	} else if f, ok := a.(bool); !ok {
		return nil, errors.New("'check_usage' must be bool")
	} else {
		policy.CheckListProjects = f
	}
	return &policy, nil
}
