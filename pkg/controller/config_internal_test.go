package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestRule_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		caseName string
		yaml     string
		exp      Rule
		isErr    bool
	}{
		{
			caseName: "normal",
			yaml: `
policy:
  type: has_projects
target: |
  github-config
`,
			exp: Rule{
				Policy: &RuleHasProjects{
					CheckListProjects: true,
					action: ActionConfig{
						Type: "fix",
					},
				},
				Target: Target{
					Patterns: []TargetPattern{
						{
							Pattern: "github-config",
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.caseName, func(t *testing.T) {
			t.Parallel()
			rule := Rule{}
			err := yaml.Unmarshal([]byte(d.yaml), &rule)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, rule)
		})
	}
}
