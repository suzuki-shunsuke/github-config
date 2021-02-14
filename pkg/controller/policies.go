package controller

type NewRepoPolicy func(param map[string]interface{}, action ActionConfig) (RepoPolicy, error)

func supportedRepoPolicies() map[string]NewRepoPolicy {
	return map[string]NewRepoPolicy{
		"has_projects": newRuleHasProjects,
	}
}
