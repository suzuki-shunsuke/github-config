package controller

import (
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/hasprojects"
)

type NewRepoPolicy func(param map[string]interface{}, action domain.ActionConfig) (domain.RepoPolicy, error)

func supportedRepoPolicies() map[string]NewRepoPolicy {
	return map[string]NewRepoPolicy{
		"has_projects": hasprojects.New,
	}
}
