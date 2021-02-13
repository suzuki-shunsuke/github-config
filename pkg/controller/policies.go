package controller

type NewRepoPolicy func() RepoPolicy

func supportedRepoPolicies() map[string]NewRepoPolicy {
	return map[string]NewRepoPolicy{
		"has_projects": func() RepoPolicy {
			return &RuleHasProjects{}
		},
	}
}
