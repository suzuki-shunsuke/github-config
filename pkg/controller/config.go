package controller

import (
	"errors"
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/github-config/pkg/domain"
	"gopkg.in/yaml.v2"
)

const fix = "fix"

type Config struct {
	Owner string `yaml:"org_name"`
	Repo  RepoConfig
	Org   OrgConfig
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	Owner          string
	GitHubToken    string
	DataDogAPIKey  string
	DryRun         bool
}

type RepoConfig struct {
	Rules []Rule
}

type Rule struct {
	Policy domain.RepoPolicy
	Target Target
}

type RuleConfig struct {
	Policy PolicyConfig
	Target Target
	Action domain.ActionConfig
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
	rule.Target = a.Target
	if a.Action.Type == "" {
		a.Action.Type = fix
	}
	newRepoMatchers := map[string]NewRepoPolicy{}
	supportedRepoPolicies(newRepoMatchers)
	if newPolicy, ok := newRepoMatchers[a.Policy.Type]; !ok {
		return errors.New("invalid policy type: " + a.Policy.Type)
	} else { //nolint:revive
		policy, err := newPolicy(a.Policy.Param, a.Action)
		if err != nil {
			return err
		}
		rule.Policy = policy
	}
	return nil
}
