package controller

import (
	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/org/defaultrepositorypermission"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/org/hasorganizationprojects"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/archived"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/hasissues"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/hasprojects"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/haswiki"
	"github.com/suzuki-shunsuke/github-config/pkg/rule/repo/private"
)

type (
	NewRepoPolicy func(param map[string]interface{}, action domain.ActionConfig) (domain.RepoPolicy, error)
	NewOrgPolicy  func(param map[string]interface{}, action domain.ActionConfig) (domain.OrgPolicy, error)
)

func supportedRepoPolicies(policies map[string]NewRepoPolicy) {
	policies["has_projects"] = hasprojects.New
	policies["has_issues"] = hasissues.New
	policies["has_wiki"] = haswiki.New
	policies["private"] = private.New
	policies["archived"] = archived.New
}

func supportedOrgPolicies(policies map[string]NewOrgPolicy) {
	policies["has_organization_projects"] = hasorganizationprojects.New
	policies["default_repository_permission"] = defaultrepositorypermission.New
}
