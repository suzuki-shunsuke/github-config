package controller

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const DataDogEventSourceType = "GITHUBCONFIG"

type RuleHasProjects struct {
	Client            *github.Client
	CheckListProjects bool
}

func (rule *RuleHasProjects) Do() []interface{} {

	return nil
}

func (rule *RuleHasProjects) doFix(updatedRepo *github.Repository) interface{} {
	return func(updatedRepo *github.Repository) error {
		updatedRepo.HasProjects = github.Bool(false)
		return nil
	}
}

func (rule *RuleHasProjects) dataDogMetrics(updatedRepo *github.Repository) interface{} {
	return nil
}

func (rule *RuleHasProjects) dataDogEvent(updatedRepo *github.Repository) *datadog.Event {
	return &datadog.Event{
		Title:      datadog.String("has_projects"),
		Text:       datadog.String("has_projects"),
		AlertType:  datadog.String("info"),
		SourceType: datadog.String(DataDogEventSourceType),
	}
}

func (rule *RuleHasProjects) match(ctx context.Context, repo *github.Repository, owner, repoName string) (bool, error) {
	if !repo.GetHasProjects() {
		return false, nil
	}
	logE := logrus.WithFields(logrus.Fields{
		"repo": repoName,
	})
	if !rule.CheckListProjects {
		return true, nil
	}
	if projects, _, err := rule.Client.Repositories.ListProjects(ctx, owner, repoName, nil); err != nil {
		logE.WithError(err).Error("list projects")
		return true, nil
	} else if len(projects) == 0 {
		return true, nil
	}
	return false, nil
}
