# github-config

[![Build Status](https://github.com/suzuki-shunsuke/github-config/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/github-config/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/github-config)](https://goreportcard.com/report/github.com/suzuki-shunsuke/github-config)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/github-config.svg)](https://github.com/suzuki-shunsuke/github-config)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/github-config/master/LICENSE)

Make GitHub Organization and Repositories Settings compliant with Policy

## Overview

`github-config` is a CLI tool to make GitHub Organization and Repositories Settings compliant with Policy.

## Install

Download binary from [Releases](https://github.com/suzuki-shunsuke/github-config/releases).

## How to use

Periodically run `github-config repo` and `github-config org`.

## Environment variables

* GITHUB_TOKEN, GITHUB_ACCESS_TOKEN
* DATADOG_API_KEY

## Usage

Please see [Usage](docs/USAGE.md)

## Action

* Fix
* Send DataDog Metrics

## Configuration

Please see [Configuration](docs/CONFIG.md)

## LICENSE

[MIT](LICENSE)
