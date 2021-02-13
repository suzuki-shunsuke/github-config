package controller

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (ctrl *Controller) RunOrg(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)))

	org, _, err := client.Organizations.Get(ctx, param.Owner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"owner": param.Owner,
		}).WithError(err).Error("get an organization")
	} else if err := ctrl.handleOrg(ctx, param, client, org); err != nil {
		logrus.WithError(err).Error("update org")
	}
	return nil
}

func (ctrl *Controller) handleOrg(ctx context.Context, param Param, client *github.Client, org *github.Organization) error { //nolint:unparam
	logE := logrus.WithFields(logrus.Fields{
		"owner": param.Owner,
	})
	updatedOrg := &github.Organization{}
	// TODO Disable projects for the organization
	// TODO Disable projects for all repositories
	if org.GetHasOrganizationProjects() {
		if projects, _, err := client.Organizations.ListProjects(ctx, param.Owner, nil); err != nil {
			logE.WithError(err).Error("list projects")
		} else if len(projects) == 0 {
			updatedOrg.HasOrganizationProjects = github.Bool(false)
		}
	}
	// TODO member privilege
	if !org.GetTwoFactorRequirementEnabled() {
		org.TwoFactorRequirementEnabled = github.Bool(true)
	}
	if org.GetDefaultRepoPermission() != "read" {
		updatedOrg.DefaultRepoPermission = github.String("read")
	}
	if org.GetMembersCanCreateRepos() {
		updatedOrg.MembersCanCreateRepos = github.Bool(false)
	}
	if org.GetMembersCanCreatePublicRepos() {
		updatedOrg.MembersCanCreatePublicRepos = github.Bool(false)
	}
	if org.GetMembersCanCreatePrivateRepos() {
		updatedOrg.MembersCanCreatePrivateRepos = github.Bool(false)
	}
	if org.GetMembersCanCreateInternalRepos() {
		updatedOrg.MembersCanCreateInternalRepos = github.Bool(false)
	}
	// TODO enable Dependabot alerts
	// TODO enable Dependabot security updates
	// TODO remove not allowed webhook url
	// TODO third party access
	// TODO github app
	// TODO Repository default branch
	// TODO Repository labels
	// TODO Team discussions
	// TODO Actions
	// TODO Organization Owned Applications
	// TODO GitHub Apps
	return nil
}
