package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const DataDogEventSourceType = "GITHUBCONFIG"

type RuleHasProjects struct {
	client  *github.Client
	datadog *datadog.Client
	// datadogMetric     datadog.Metric
	// datadogEvent      datadog.Event
	Actions                []ActionConfig
	CheckListProjects      bool
	hasActionDataDogMetric bool
}

func (rule *RuleHasProjects) SetGitHubClient(client *github.Client) {
	rule.client = client
}

func (rule *RuleHasProjects) SetDataDogClient(client *datadog.Client) {
	rule.datadog = client
}

func (rule *RuleHasProjects) DataDogMetric(ctx context.Context, param *ParamAction) error {
	if rule.hasActionDataDogMetric {
		param.DataDogMetrics = append(param.DataDogMetrics, rule.dataDogMetric(param.Repo, &param.TimestampFloat64))
	}
	return nil
}

func (rule *RuleHasProjects) Action(ctx context.Context, param *ParamAction) error {
	logE := logrus.WithFields(logrus.Fields{
		"repo":   param.Repo.Name,
		"org":    param.Repo.Owner,
		"policy": "has_projects",
	})
	for _, action := range rule.Actions {
		switch action.Type {
		case "datadog_metric":
		case "fix":
			rule.Fix(ctx, param)
		case "datadog_event":
			if param.DryRun {
				logE.Info("[DRY RUN] post an event to DataDog")
			} else {
				ev := rule.DataDogEvent(param.Repo, param.TimestampInt)
				if _, err := rule.datadog.PostEvent(&ev); err != nil {
					return fmt.Errorf("post an event to DataDog: %w", err)
				}
				logE.Info("post an event to DataDog")
			}
		default:
			return errors.New("invalid action type: " + action.Type)
		}
	}
	return nil
}

func (rule *RuleHasProjects) Fix(ctx context.Context, param *ParamAction) {
	param.UpdatedRepo.HasProjects = github.Bool(false)
	param.IsEdited = true
}

func (rule *RuleHasProjects) DataDogEvent(repo Repository, eventTime int) datadog.Event {
	return datadog.Event{
		Title: datadog.String("github-config rule violation is detected. GitHub Projects is enabled"),
		Text: datadog.String(`Disable GitHub Projects if it isn't needed.
If this is false positive, please fix github-config's configuration.
https://github.com/suzuki-shunsuke/github-config`),
		Time:       &eventTime,
		Priority:   datadog.String("low"),
		AlertType:  datadog.String("warning"),
		SourceType: datadog.String(DataDogEventSourceType),
		Tags: []string{
			"github_repo:" + repo.Name,
			"github_org:" + repo.Owner,
		},
	}
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

func newRuleHasProjects(param map[string]interface{}, actions []ActionConfig) (RepoPolicy, error) {
	policy := RuleHasProjects{
		Actions:                actions,
		hasActionDataDogMetric: hasDataDogMetricAction(actions),
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
