package domain

import (
	"context"

	"github.com/google/go-github/v33/github"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

type ActionConfig struct {
	Type  string
	Param map[string]interface{}
}

type RepoPolicy interface {
	Match(ctx context.Context, repo Repository) (bool, error)
	Action() ActionConfig
	DataDogMetric(repo Repository, now *float64) datadog.Metric
}

type Fixable interface {
	Fix(ctx context.Context, param *ParamAction)
}

type UseGitHubClient interface {
	SetGitHubClient(*github.Client)
}

type UseDataDogClient interface {
	SetDataDogClient(*datadog.Client)
}

type Repository struct {
	GitHub *github.Repository
	Owner  string
	Name   string
}

type ParamAction struct {
	Repo             Repository
	UpdatedRepo      *github.Repository
	Actions          []ActionConfig
	DataDogMetrics   []datadog.Metric
	TimestampFloat64 float64
	TimestampInt     int
	IsEdited         bool
	DryRun           bool
}
