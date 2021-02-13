package controller

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/v33/github"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Owner string `yaml:"org_name"`
	Repo  RepoConfig
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	Owner          string
	GitHubToken    string
	DryRun         bool
}

type RepoConfig struct {
	Rules []Rule
}

type RepoPolicy interface {
	Match(ctx context.Context, repo Repository) (bool, error)
	SetGitHubClient(*github.Client)
}

type Rule struct {
	Policy  RepoPolicy
	Targets Targets
}

type RuleConfig struct {
	Policy  PolicyConfig
	Targets Targets
}

type PolicyConfig struct {
	Type  string
	Param map[string]interface{}
}

func (ctrl *Controller) readConfig(param Param, cfg *Config) error {
	cfgFile, err := os.Open(param.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("open a configuration file %s: %w", param.ConfigFilePath, err)
	}
	defer cfgFile.Close()
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return fmt.Errorf("parse a configuration file as YAML %s: %w", param.ConfigFilePath, err)
	}
	return nil
}

func (rule *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	a := RuleConfig{}
	if err := unmarshal(&a); err != nil {
		return err
	}
	rule.Targets = a.Targets
	newRepoMatchers := supportedRepoPolicies()
	if newPolicy, ok := newRepoMatchers[a.Policy.Type]; !ok {
		return errors.New("invalid policy type: " + a.Policy.Type)
	} else {
		policy, err := newPolicy(a.Policy.Param)
		if err != nil {
			return err
		}
		rule.Policy = policy
	}
	return nil
}
