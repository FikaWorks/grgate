<p align=center><img src="https://github.com/fikaworks/grgate/raw/main/.github/grgate-logo-300x300.png" alt="GRGate, git release gate" width="300" height="300"></p>

GRGate - Git release gate
=========================

> **grgate**, a git release gate utility which autopublish draft/unpublished
releases based on commit status (aka checks). It can be triggered automatically
using Git webhook or directly from the CLI.

Currently, only Github and Gitlab are supported, other provider could come in a
near future.

![grgate workflow](https://github.com/FikaWorks/grgate/actions/workflows/main.yml/badge.svg?branch=main)

## Overview

GRGate is a CLI which can run a server and listen to Git webhook. When a release
is published as draft, GRGate will wait for all the status check attached to the
commit target of the release to succeed before merging it.

The following diagram represent a concret example where a CI/CD process
generate/publish versionned artifacts and generate a draft release. Artifacts
are then deployed by a third party to different environments running with
different config. End-to-end tests are then run against these 2 environments
and reports result to the draft release as commit status. When all tests pass,
GRGate publish the Github release.

![GRGate Overview](grgate-overview.png)

### Unpublished releases terminology

Different terminology is used by different provider:

- **Github** uses the term [draft releases][draft-release] to prepare a release
without publishing it.
- **Gitlab** uses the term [upcoming releases][upcoming-release], it is similar
to Github Pre-releases where a badge notify the upcoming release in the Gitlab
release page.  The attribute `released_at` should be set to a future date to
have it enabled and it is only possible to change it using the Gitlab API.

[draft-release]: https://docs.github.com/en/github/administering-a-repository/managing-releases-in-a-repository#about-release-management
[upcoming-release]: https://docs.gitlab.com/ee/api/releases/#upcoming-releases

## Getting started

Create a Github APP or Gitlab token, then create a `config.yaml` or a
`values.yaml` file if using Helm charts.

### Helm chart

A Helm chart is available in the [FikaWorks Helm charts
repository][helm-charts].

[helm-charts]: https://github.com/FikaWorks/helm-charts

```bash
$ helm repo add fikaworks https://fikaworks.github.io/helm-charts
$ helm install --name grgate --f my-values.yaml fikaworks/grgate
```

### GRGate CLI

Download latest release from the [release page][release-page].

[release-page]: https://github.com/fikaworks/grgate/releases

```bash
# check available commands
$ grgate --help

# run GRGate against the Github repository FikaWorks/grgate:
$ grgate run FikaWorks/grgate

# listen to Git webhook events on 0.0.0.0:8080
$ grgate serve -c config.yaml

# list status for a given commit in the FikaWorks/grgate repository
$ grgate status list FikaWorks/grgate \
    --commit 93431f42d5a5abc2bb6703fc723b162a9d2f20c3

# set status of given commit (Github)
$ grgate status set FikaWorks/grgate \
    --commit 93431f42d5a5abc2bb6703fc723b162a9d2f20c3 \
    --name e2e-happyflow \
    --status completed \
    --state success
```

## Config reference

GRGate has 2 types of configuration:
- **main config** which define global, server settings and credentials to talk
to external API like Github
- **repo config** which is a file stored in the repository and define the
merging rule. Fields are inherited from `main.globals:`, so if you need to
override a setting you can define it in the repository itself

### Main config

The main configuration file can be passed to the CLI via the `--config`
argument, by default it will try to read from `/etc/grgate/config.yaml`.

```yaml
# global configuration, this is the default
globals:
  # enable GRGatem, if set to false release are not published
  enabled: true

  # filter release by tag, the tag associated to the draft/unpublished releases
  # must match the regular expression defined by tagRegexp, default: .*
  tagRegexp: .*

  # list of statuses, default: []
  statuses:
    - e2e happy flow

  # append statuses to release note
  releaseNote:
    enabled: true
    template: |-
      {{- .ReleaseNote -}}
      <!-- GRGate start -->
      <details><summary>Status check</summary>
      {{ range .Statuses }}
      - [{{ if or (eq .Status "completed" ) (eq .Status "success") }}x{{ else }} {{ end }}] {{ .Name }}
      {{- end }}

      </details>
      <!-- GRGate end -->

# server configuration (webhook)
# webhook should be sent to /<provider>/webhook, where provider is either
# github or gitlab
server:
  listenAddress: 0.0.0.0:8080
  metricsAddress: 0.0.0.0:9101
  probeAddress: 0.0.0.0:8086
  webhookSecret: a-random-string

# number of workers to run, default: 1
workers: 1

# platform to use
platform: github # github|gitlab, default: github

# Github configuration
# when creating the Github app, make sure to set the following permissions:
#   - Administration read/write
#   - Checks read/write
#   - Contents read/write
#   - Metadata read-only
#   - Commit statuses read/write
#
# subscribe to the following webhook events:
#   - Check runs
#   - Check suites
#   - Releases
#   - Statuses
github:
  appID: 000000
  installationID: 00000000
  privateKeyPath: path-to-key.pem

# Gitlab configuration
# when creating the Gitlab token, make sure to set the following permissions:
#   - read_repository
# subscribe to the following webhook events:
#   - Release events
#   - Pipeline events
gitlab:
  token: gitlab-token

# configuration can be overriden in the repository itself, you can define the
# default path below, default: .grgate.yaml
repoConfigPath: .grgate.yaml

logLevel: info  # trace|debug|info|warn|error|fatal|panic, default: info
logFormat: json # json|pretty, default: pretty
```

### Repository config

The `globals:` section of the GRGate configuration can be overriden from the
repository itself. If you create a file named `.grgate.yaml` at the root of the
repository, GRGate will read it before processing the repository.

```yaml
# only process releases with tag matching a regular expression pattern
tagRegexp: v1.0.0-beta-\d*

# automerge release if the following status succeeded
statuses:
  - e2e-with-feature-A-on
  - e2e-with-feature-B-on
```

## Contributing / local development

For local development and to contribute to this project, refer to
[CONTRIBUTING.md](./CONTRIBUTING.md).
