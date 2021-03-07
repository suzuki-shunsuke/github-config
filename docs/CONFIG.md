# Configuration

github-config.yaml

ex.

```yaml
---
org_name: terraform-provider-graylog
org:
  rules:
  - policy:
      type: has_organization_projects
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
    target: |
      test-circleci
    action:
      type: datadog_metric
```
