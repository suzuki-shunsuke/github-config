package controller

type Targets []string

func (targets Targets) Match(repo Repository) (bool, error) {
	for _, target := range targets {
		if repo.Name == target {
			return true, nil
		}
	}
	return false, nil
}
