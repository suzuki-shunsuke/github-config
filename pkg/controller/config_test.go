package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/github-config/pkg/controller"
	"gopkg.in/yaml.v2"
)

func TestRule_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		caseName string
		yaml     string
		exp      controller.Rule
		isErr    bool
	}{
		{
			caseName: "normal",
			yaml: `
policy:
  type: has_projects
targets:
- github-config
`,
			exp: controller.Rule{
				Policy:  &controller.RuleHasProjects{},
				Targets: controller.Targets{"github-config"},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.caseName, func(t *testing.T) {
			t.Parallel()
			rule := controller.Rule{}
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
