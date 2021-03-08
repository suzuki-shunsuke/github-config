package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"golang.org/x/oauth2"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const datadogMetric = "datadog_metric"

func (ctrl *Controller) RunOrg(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"count": len(cfg.Org.Rules),
	}).Info("list org rules")

	ctrl.Config = cfg
	param.Owner = cfg.Owner

	if param.DataDogAPIKey != "" {
		ctrl.DataDog = datadog.NewClient(param.DataDogAPIKey, "")
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: param.GitHubToken},
	)))
	for i, rule := range cfg.Org.Rules {
		if policy, ok := rule.Policy.(domain.UseGitHubClient); ok {
			policy.SetGitHubClient(client)
		}
		if policy, ok := rule.Policy.(domain.UseDataDogClient); ok {
			policy.SetDataDogClient(ctrl.DataDog)
		}
		cfg.Org.Rules[i] = rule
	}

	org, err := ctrl.getOrg(ctx, client, param)
	if err != nil {
		return err
	}
	if err := ctrl.handleOrg(ctx, param, client, org); err != nil {
		return fmt.Errorf("handle an organization: %w", err)
	}
	return nil
}

func (ctrl *Controller) getOrg(ctx context.Context, client *github.Client, param Param) (domain.Organization, error) {
	org, _, err := client.Organizations.Get(ctx, param.Owner)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("get an organization (owner: %s): %w", param.Owner, err)
	}
	return domain.Organization{
		GitHub: org,
		Name:   param.Owner,
	}, nil
}

func (ctrl *Controller) orgAction(ctx context.Context, param *domain.ParamOrgAction, policy domain.OrgPolicy) error {
	switch t := policy.Action().Type; t {
	case datadogMetric:
	case "fix":
		a, ok := policy.(domain.OrgFixable)
		if !ok {
			return errors.New("this rule doesn't support to fix")
		}
		a.Fix(ctx, param)
	default:
		return errors.New("invalid action type: " + t)
	}
	return nil
}

func (ctrl *Controller) handleOrg(ctx context.Context, param Param, client *github.Client, org domain.Organization) error { //nolint:unparam,cyclop
	orgName := org.Name
	logE := logrus.WithFields(logrus.Fields{
		"org": orgName,
	})
	ts := time.Now().Unix()
	paramAction := domain.ParamOrgAction{
		Org:              org,
		UpdatedOrg:       &github.Organization{},
		TimestampFloat64: float64(ts),
		TimestampInt:     int(ts),
		DryRun:           param.DryRun,
	}
	for _, rule := range ctrl.Config.Org.Rules {
		actionConfig := rule.Policy.Action()
		logE.Debug("check rule")
		if actionConfig.Type == datadogMetric {
			paramAction.DataDogMetrics = append(paramAction.DataDogMetrics, rule.Policy.DataDogMetric(paramAction.Org, &paramAction.TimestampFloat64))
		}
		if f, err := rule.Policy.Match(ctx, org); err != nil {
			logE.WithError(err).Error("check an organization matches with the policy")
			continue
		} else if f {
			logE.Info("an organization matches with the rule")
			if err := ctrl.orgAction(ctx, &paramAction, rule.Policy); err != nil {
				logE.WithError(err).Error("prepare")
				continue
			}
		}
	}
	if paramAction.IsEdited {
		if param.DryRun {
			logE.Info("[DRY RUN] update an organization")
		} else {
			if _, _, err := client.Organizations.Edit(ctx, orgName, paramAction.UpdatedOrg); err != nil {
				logE.WithError(err).Error("update an organization")
			}
			logE.Info("update an organization")
		}
	}
	if len(paramAction.DataDogMetrics) != 0 {
		if param.DryRun {
			logE.Info("[DRY RUN] post metrics to DataDog")
		} else {
			if err := ctrl.DataDog.PostMetrics(paramAction.DataDogMetrics); err != nil {
				logE.WithError(err).Error("post metrics to DataDog")
			}
			logE.Info("post metrics to DataDog")
		}
	}
	return nil
}
