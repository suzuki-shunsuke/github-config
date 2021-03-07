package controller

import (
	"errors"

	"github.com/suzuki-shunsuke/github-config/pkg/domain"
)

type OrgConfig struct {
	Rules []OrgRule
}

type OrgRule struct {
	Policy domain.OrgPolicy
}

func (rule *OrgRule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	a := RuleConfig{}
	if err := unmarshal(&a); err != nil {
		return err
	}
	if a.Action.Type == "" {
		a.Action.Type = fix
	}
	newOrgMatchers := map[string]NewOrgPolicy{}
	supportedOrgPolicies(newOrgMatchers)
	if newPolicy, ok := newOrgMatchers[a.Policy.Type]; !ok {
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
