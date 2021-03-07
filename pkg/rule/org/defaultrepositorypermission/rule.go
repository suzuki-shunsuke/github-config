package defaultrepositorypermission

import (
	"context"
	"errors"

	"github.com/google/go-github/v33/github"
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const datadogMetricName = "github_config.org.default_repository_permission"

type Rule struct {
	action     domain.ActionConfig
	permission string
}

func New(param map[string]interface{}, action domain.ActionConfig) (domain.OrgPolicy, error) {
	policy := Rule{
		action: action,
	}
	if a, ok := param["permission"]; !ok {
		policy.permission = "read"
		return &policy, nil
	} else if f, ok := a.(string); !ok {
		return nil, errors.New("'permission' must be string")
	} else {
		policy.permission = f
	}
	return &policy, nil
}

func (rule *Rule) DataDogMetric(org domain.Organization, now *float64) datadog.Metric {
	f := 0.0
	if rule.permission != org.GitHub.GetDefaultRepoPermission() {
		f = 1.0
	}
	return datadog.Metric{
		Metric: datadog.String(datadogMetricName),
		Points: []datadog.DataPoint{
			{now, &f},
		},
		Tags: []string{
			"github_org:" + org.Name,
			"permission:" + org.GitHub.GetDefaultRepoPermission(),
		},
	}
}

func (rule *Rule) Action() domain.ActionConfig {
	return rule.action
}

func (rule *Rule) Fix(ctx context.Context, param *domain.ParamOrgAction) {
	param.UpdatedOrg.DefaultRepoPermission = github.String(rule.permission)
	param.IsEdited = true
}

func (rule *Rule) Match(ctx context.Context, org domain.Organization) (bool, error) {
	return org.GitHub.GetDefaultRepoPermission() != rule.permission, nil
}
