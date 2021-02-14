package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v33/github"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

type RuleHasProjects struct {
	client  *github.Client
	datadog *datadog.Client
	// datadogMetric     datadog.Metric
	action            ActionConfig
	CheckListProjects bool
}

func newRuleHasProjects(param map[string]interface{}, action ActionConfig) (RepoPolicy, error) {
	policy := RuleHasProjects{
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

func (rule *RuleHasProjects) SetGitHubClient(client *github.Client) {
	rule.client = client
}

func (rule *RuleHasProjects) SetDataDogClient(client *datadog.Client) {
	rule.datadog = client
}

func (rule *RuleHasProjects) DataDogMetric(ctx context.Context, param *ParamAction) error {
	if rule.action.Type == "datadog_metric" {
		param.DataDogMetrics = append(param.DataDogMetrics, rule.dataDogMetric(param.Repo, &param.TimestampFloat64))
	}
	return nil
}

func (rule *RuleHasProjects) dataDogMetric(repo Repository, now *float64) datadog.Metric {
	f := 0.0
	if repo.GitHub.GetHasProjects() {
		f = 1.0
	}
	return datadog.Metric{
		Metric: datadog.String("github_config.repo.has_projects"),
		Points: []datadog.DataPoint{
			{now, &f},
		},
		Tags: []string{
			"github_repo:" + repo.Name,
			"github_org:" + repo.Owner,
		},
	}
}

func (rule *RuleHasProjects) Action(ctx context.Context, param *ParamAction) error {
	switch rule.action.Type {
	case "datadog_metric":
	case "fix":
		rule.Fix(ctx, param)
	default:
		return errors.New("invalid action type: " + rule.action.Type)
	}
	return nil
}

func (rule *RuleHasProjects) Fix(ctx context.Context, param *ParamAction) {
	param.UpdatedRepo.HasProjects = github.Bool(false)
	param.IsEdited = true
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
