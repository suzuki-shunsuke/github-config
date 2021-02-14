package controller

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Target struct {
	Patterns []TargetPattern
}

func (target *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s := ""
	if err := unmarshal(&s); err != nil {
		return err
	}
	var patterns []TargetPattern //nolint:prealloc
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "!") {
			patterns = append(patterns, TargetPattern{
				Not:     true,
				Pattern: line[1:],
			})
			continue
		}
		patterns = append(patterns, TargetPattern{
			Pattern: line,
		})
	}
	target.Patterns = patterns
	return nil
}

func (target *Target) Match(repo Repository) (bool, error) {
	matched := false
	for _, pattern := range target.Patterns {
		if matched {
			if !pattern.Not {
				continue
			}
			if f, err := filepath.Match(pattern.Pattern, repo.Name); err != nil {
				return false, fmt.Errorf("check whether the repository name matched with pattern: %w", err)
			} else if f {
				matched = false
				continue
			}
			continue
		}
		if pattern.Not {
			continue
		}
		if f, err := filepath.Match(pattern.Pattern, repo.Name); err != nil {
			return false, fmt.Errorf("check whether the repository name matched with pattern: %w", err)
		} else if f {
			matched = true
			continue
		}
	}
	return matched, nil
}

type TargetPattern struct {
	Not     bool
	Pattern string
}

func (target *TargetPattern) Match(repo Repository) (bool, error) {
	if f, err := filepath.Match(target.Pattern, repo.Name); err != nil {
		return false, fmt.Errorf("check whether repsitory name matches with pattern: %w", err)
	} else if f {
		return !target.Not, nil
	}
	return target.Not, nil
}
