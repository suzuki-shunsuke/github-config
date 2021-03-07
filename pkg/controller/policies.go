package controller

import (
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/archived"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/hasissues"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/hasprojects"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/haswiki"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/private"
)

type NewRepoPolicy func(param map[string]interface{}, action domain.ActionConfig) (domain.RepoPolicy, error)

func supportedRepoPolicies() map[string]NewRepoPolicy {
	return map[string]NewRepoPolicy{
		"has_projects": hasprojects.New,
		"has_issues":   hasissues.New,
		"has_wiki":     haswiki.New,
		"private":      private.New,
		"archived":     archived.New,
	}
}
