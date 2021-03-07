package haswiki

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const datadogMetricName = "github_config.repo.has_wiki"

type Rule struct {
	action domain.ActionConfig
}

func New(param map[string]interface{}, action domain.ActionConfig) (domain.RepoPolicy, error) {
	policy := Rule{
		action: action,
	}
	return &policy, nil
}

func (rule *Rule) DataDogMetric(repo domain.Repository, now *float64) datadog.Metric {
	f := 0.0
	if repo.GitHub.GetHasWiki() {
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
	param.UpdatedRepo.HasWiki = github.Bool(false)
	param.IsEdited = true
}

func (rule *Rule) Match(ctx context.Context, repo domain.Repository) (bool, error) {
	return repo.GitHub.GetHasWiki(), nil
}
