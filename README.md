GGate - Git release gate
========================

> **ggate**, is git release gate utility which autopublish draft/unpublished
releases based on commit status (aka checks). It can be triggered automatically
using Git webhook or directly from the CLI.

Currently, only the Github platform is supported, Gitlab and other platform
will in a near future.

## Overview

GGate is a CLI which can run a server and listen to Git webhook. When a release
is published as draft, GGate will wait for all the status check attached to the
the commit target of the release to succeed before merging it.

![GGate Overview](ggate-overview.png)

## Getting started

Download latest release from the [release page][0].

[0]: https://github.com/fikaworks/ggate/releases

```bash
# check available commands
$ ggate --help

# run GGate against the Github repository FikaWorks/ggate:
$ ggate run FikaWorks/ggate

# listen to Git webhook events on 0.0.0.0:8080
$ ggate serve

# list status for a given commit in the FikaWorks/ggate repository
$ ggate status list FikaWorks/ggate \
    --commit 93431f42d5a5abc2bb6703fc723b162a9d2f20c3

# set status of given commit
$ ggate status set FikaWorks/ggate \
    --commit 93431f42d5a5abc2bb6703fc723b162a9d2f20c3 \
    --name e2e-happyflow \
    --status completed \
    --state success
```

## Config reference

### Local CLI / Server

The configuration file can be passed to the CLI via the `--config` argument, by
default it will try to read `~/.ggate.yaml` from the home directory.

```yaml
# global configuration, this is the default
globals:
  # filter release by tag, the tag associated to the draft/unpublished releases
  # must match the regular expression defined by tagRegexp, default: .*
  tagRegexp: .*

  # list of statuses, default: []
  statuses:
    - e2e happy flow

server:
  listenAddress: 0.0.0.0:8080
  metricsAddress: 0.0.0.0:9101
  probeAddress: 0.0.0.0:8086

# number of workers to run, default: 1
workers: 1

# Github configuration
# when creating the Github app, make sure to select the following webhook
# events:
#   - Check runs
#   - Check suites
#   - Releases
#   - Statuses
github:
  appID: 000000
  installationID: 00000000
  privateKeyPath: path-to-key.pem
  webhookSecret: a-random-string

# configuration can be overriden in the repository itself, you can define the
# default path below, default: .ggate.yaml
repoConfigPath: .ggate.yaml

logLevel: info  # trace|debug|info|warn|error|fatal|panic, default: info
logFormat: json # json|pretty, default: pretty
```

### Repository

The `globals:` section of the GGate configuration can be overriden from the
repository itself. If you create a file named `.ggate.yaml` at the root of the
repository, GGate will read it before processing the repository.

```yaml
# only process releases with tag matching a regular expression pattern
tagRegexp: v1.0.0-beta-\d*

# automerge release if the following status succeeded
statuses:
  - e2e-with-feature-A-on
  - e2e-with-feature-B-on
```

## Development

For local development, you can use [ngrok](https://ngrok.com/) to receive
Github webhook events.

Start ngrok and forward requests to port `8080`:

```bash
$ ngrok http 8080
```

Create a new webhook from the Github organization or repository settings page.

Use the following settings:
- set the payload URL to your ngrok endpoint, ie:
  `http://bae0d008e18b.ngrok.io/github/webhook`
- content type: `application/json`
- select individual events:
  - Check runs
  - Check suites
  - Releases
  - Statuses

Run the server locally:
```bash
$ go run main.go serve
```

If you create/update a status check or create a draft release, you should see
GGate processing the triggering repository.

### Build binary

```bash
$ make build
$ ./ggate --help
```

### Docker

```bash
$ make build-docker
$ docker run --ti -p 8080:8080 fikaworks/ggate
```

### Tests, lint, vet

```bash
$ make validate
```
