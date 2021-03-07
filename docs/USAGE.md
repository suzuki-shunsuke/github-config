# Usage

```
$ github-config help
NAME:
   github-config - Make GitHub Organization and Repositories Settings compliant with Policy. https://github.com/suzuki-shunsuke/github-config

USAGE:
   github-config [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   repo     Make GitHub Repositories Settings compliant with Policy
   org      Make GitHub Organization Settings compliant with Policy
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## org

```
$ github-config help org
NAME:
   github-config org - Make GitHub Organization Settings compliant with Policy

USAGE:
   github-config org [command options] [arguments...]

OPTIONS:
   --log-level value         log level
   --config value, -c value  configuration file path
   --dry-run                 dry run (default: false)
```

## repo

```
$ github-config help repo
NAME:
   github-config repo - Make GitHub Repositories Settings compliant with Policy

USAGE:
   github-config repo [command options] [arguments...]

OPTIONS:
   --log-level value         log level
   --config value, -c value  configuration file path
   --dry-run                 dry run (default: false)
```
