# Configuration

## Configuration file path

By default, `github-config` reads `github-config.yaml` on the current directory.
We can specify the path by the command line option `--config (-c)`.

## Repository Policy

[Repository Policy](CONFIG_REPO_POLICY.md)

## Organization Policy

[Organization Policy](CONFIG_ORG_POLICY.md)

## Example

```yaml
---
org_name: terraform-provider-graylog
org:
  rules:
  - policy:
      type: has_organization_projects
      param:
        check_usage: true
    action:
      type: fix
  - policy:
      type: default_repository_permission
    action:
      type: fix
repo:
  rules:
  - policy:
      type: has_projects
      param:
        check_usage: true
    target: |
      terraform-provider-graylog
    action:
      type: datadog_metric
```

## Reference

path | type | required | default | description
--- | --- | --- | --- | ---
.org_name | string | true | organizaiton name
.org | []Org | false | [] | organization config
.repo | []Repo | false | [] | repository config

## Type: Org

path | type | required | default | description
--- | --- | --- | --- | ---
.rules | []OrgRule | false | [] | organization rules

## Type: OrgRule

path | type | required | default | description
--- | --- | --- | --- | ---
.policy | OrgPolicy | true | | organization policy
.action | Action | true | | action

## Type: OrgPolicy

path | type | required | default | description
--- | --- | --- | --- | ---
.type | string | true | | policy type
.param | map[string]interface{} | false | {} | policy parameter. The type depends on the policy

Please see [Organization Policy](CONFIG_ORG_POLICY.md) too.

## Type: Action

path | type | required | default | description
--- | --- | --- | --- | ---
.type | string | false | "fix" | action type ("fix" or "datadog_metric")

* `fix`: the setting is fixed automatically
* `datadog_metric`: the metrics is sent to DataDog

## Type: Repo

path | type | required | default | description
--- | --- | --- | --- | ---
.rules | []RepoRule | false | [] | repository rules

## Type: RepoRule

path | type | required | default | description
--- | --- | --- | --- | ---
.policy | RepoPolicy | true | | repository policy
.action | Action | true | | action
.target | string | true | | expression to filter target repositories

## Type: RepoPolicy

path | type | required | default | description
--- | --- | --- | --- | ---
.type | string | true | | policy type
.param | map[string]interface{} | false | {} | policy parameter. The type depends on the policy

Please see [Repository Policy](CONFIG_REPO_POLICY.md) too.

## RepoRule.target - expression to filter target repositories

ex.

```yaml
    target: |
      *
      !test-*
```

RepoRule.target is an expression to filter target repositories.
The expression is similar to [gitignore](https://github.com/git/git/blob/v2.19.1/Documentation/gitignore.txt#L70).

* line starts with `#` is ignored
* `!` any matching repository included by a previous pattern will become excluded again
* [filepath.Match](https://golang.org/pkg/path/filepath/#Match) is used to judge whether the repository name matches with the pattern
* if repository name matches with any pattern, the repository is proceeded by the policy
