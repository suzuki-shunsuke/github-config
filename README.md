# github-config

GitHub Organization Configuration management tool

## Overview

## Install

Download binary from [Releases]().

## Required permission of GitHub Token

## DataDog Events

https://docs.datadoghq.com/events/

## How to use

Periodically run `github-config repo` and `github-config org`.

## Action

* Fix
* Notify

## Environment variables

* GITHUB_TOKEN, GITHUB_ACCESS_TOKEN
* DATADOG_API_KEY

## Configuration

github-config.yaml

```yaml
---
org_name: suzuki-shunsuke
datadog: true
org:
  items:
  - rule: has_organization_projects
    enabled: true
    action:
      fix: true
  - rule:
      type: two_factor_requirement_enabled
    enabled: true
    action:
      fix: false
      datadog_event:
        enabled: true
        param:
          # https://docs.datadoghq.com/api/latest/events/#post-an-event
          aggregation_key:
          alert_type:
          priority:
          source_type_name:
          tags:
          - "owner:sre"
          text:
          title:
repo:
  items:
  - rule:
      type: visibility
      param:
        visibility: private
    condition:
      exclude: ["foo"]
    action:
```

## LICENSE

[MIT](LICENSE)
