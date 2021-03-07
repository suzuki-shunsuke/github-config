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

type OrgPolicy interface {
	Match(ctx context.Context, org Organization) (bool, error)
	Action() ActionConfig
	DataDogMetric(org Organization, now *float64) datadog.Metric
}

type Fixable interface {
	Fix(ctx context.Context, param *ParamAction)
}

type OrgFixable interface {
	Fix(ctx context.Context, param *ParamOrgAction)
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

type Organization struct {
	GitHub *github.Organization
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

type ParamOrgAction struct {
	Org              Organization
	UpdatedOrg       *github.Organization
	Actions          []ActionConfig
	DataDogMetrics   []datadog.Metric
	TimestampFloat64 float64
	TimestampInt     int
	IsEdited         bool
	DryRun           bool
}
