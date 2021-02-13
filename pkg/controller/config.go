package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GitHubToken string
	Owner       string `yaml:"org_name"`
	Org         Org
	Repo        Repo
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	Owner          string
}

type Org struct {
	Items []OrgItem
}

type OrgItem struct {
	Rule    RuleW
	Action  Action
	Enabled bool
}

type Action struct {
	Fix          bool
	DataDogEvent DataDogEvent `yaml:"datadog_event"`
}

type DataDogEvent struct {
	Enabled bool
}

type DataDogEventParam struct {
	AggregationKey string `yaml:"aggregation_key"`
	AlertType      string `yaml:"alert_type"`
	Priority       string `yaml:"priority"`
	SourceTypeName string `yaml:"source_type_name"`
	Tags           []string
	Text           string
	Title          string
}

type Repo struct {
	Items []RepoItem
}

type RepoItem struct {
	Rule      RuleW
	Condition Condition
	Action    Action
}

func (ctrl *Controller) readConfig(param Param, cfg *Config) error {
	cfgFile, err := os.Open(param.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("open a configuration file %s: %w", param.ConfigFilePath, err)
	}
	defer cfgFile.Close()
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return fmt.Errorf("parse a configuration file as YAML %s: %w", param.ConfigFilePath, err)
	}
	return nil
}

type Condition struct {
	Exclude []string
	Include []string
}

func (cond *Condition) MatchRepo(repo *github.Repository) (bool, error) {
	repoName := repo.GetName()
	if len(cond.Exclude) > 0 {
		for _, exclude := range cond.Exclude {
			if repoName == exclude {
				return false, nil
			}
		}
		return true, nil
	}
	if len(cond.Include) > 0 {
		for _, include := range cond.Include {
			if repoName == include {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}

type RuleW struct {
	Type  string
	Param RuleParam
}

type RuleParam map[string]interface{}

func (rule *RuleW) MatchRepo(ctx context.Context, repo *github.Repository, client *github.Client) (bool, error) {
	owner := repo.GetOwner().GetName()
	repoName := repo.GetName()
	logE := logrus.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repoName,
	})
	switch rule.Type {
	case "has_projects":
		if repo.GetHasProjects() {
			if projects, _, err := client.Repositories.ListProjects(ctx, owner, repoName, nil); err != nil {
				logE.WithError(err).Error("list projects")
				return true, nil
			} else if len(projects) == 0 {
				return true, nil
			}
		}
	case "has_issues":
		if repo.GetHasIssues() {
			if issues, _, err := client.Issues.ListByRepo(ctx, owner, repoName, nil); err != nil {
				logE.WithError(err).Error("list issues")
				return true, nil
			} else if len(issues) == 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func (repoItem *RepoItem) DoAction(ctx context.Context, repo, updatedRepo *github.Repository, client *github.Client) error {
	owner := repo.GetOwner().GetName()
	repoName := repo.GetName()
	logE := logrus.WithFields(logrus.Fields{
		"owner": owner,
		"repo":  repoName,
	})
	switch repoItem.Rule.Type {
	case "has_projects":
		if repo.GetHasProjects() {
			if projects, _, err := client.Repositories.ListProjects(ctx, owner, repoName, nil); err != nil {
				logE.WithError(err).Error("list projects")
				return nil
			} else if len(projects) == 0 {
				updatedRepo.HasProjects = github.Bool(false)
				return nil
			}
		}
	}
	return nil
}
