package hasprojects

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const datadogMetricName = "github_config.repo.has_projects"

type Rule struct {
	client            *github.Client
	action            domain.ActionConfig
	CheckListProjects bool
}

func New(param map[string]interface{}, action domain.ActionConfig) (domain.RepoPolicy, error) {
	policy := Rule{
		action: action,
	}
	if a, ok := param["check_usage"]; !ok {
		policy.CheckListProjects = true
		return &policy, nil
	} else if f, ok := a.(bool); !ok {
		return nil, errors.New("'check_usage' must be bool")
	} else {
		policy.CheckListProjects = f
	}
	return &policy, nil
}

func (rule *Rule) SetGitHubClient(client *github.Client) {
	rule.client = client
}

func (rule *Rule) DataDogMetric(repo domain.Repository, now *float64) datadog.Metric {
	f := 0.0
	if repo.GitHub.GetHasProjects() {
		f = 1.0
	}
	return datadog.Metric{
		Metric: datadog.String(datadogMetricName),
		Points: []datadog.DataPoint{
			{now, &f},
		},
		Tags: []string{
			"github_repo:" + repo.Name,
			"github_org:" + repo.Owner,
		},
	}
}

func (rule *Rule) Action() domain.ActionConfig {
	return rule.action
}

func (rule *Rule) Fix(ctx context.Context, param *domain.ParamAction) {
	param.UpdatedRepo.HasProjects = github.Bool(false)
	param.IsEdited = true
}

func (rule *Rule) Match(ctx context.Context, repo domain.Repository) (bool, error) {
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
