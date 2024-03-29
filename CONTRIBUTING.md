Contributing to GRGate
======================

### Local development

For local development, you can use [ngrok](https://ngrok.com/) to receive
GitHub webhook events.

Start ngrok and forward requests to port `8080`:

```bash
$ ngrok http 8080
```

Create a new webhook from the GitHub organization or repository settings page.

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
GRGate processing the triggering repository.

### Build binary

```bash
$ make build
$ ./grgate --help
```

### Docker

```bash
$ make build-docker
$ docker run --ti -p 8080:8080 -v $PWD/config.yaml:/etc/grgate/config.yaml ghcr.io/fikaworks/grgate
```

### Lint

[golangci-lint](https://golangci-lint.run) is used to lint the project, make
sure to have the binary installed on your machine.

```bash
$ make lint
```

You can let `golangci-lint` automatically fix linting issue using the following
command:

```bash
$ golangci-lint run --fix
```

### Tests

Run unit tests:

```bash
$ make test
```

[gomock](https://github.com/golang/mock) mocking framework is used to generate
mocks, you can regenerate all mocks by using the following command:

```bash
$ go install github.com/golang/mock/mockgen@latest
$ make mocks
```

### Integration tests

Run integration tests against all platforms. The tests create temporary
repositories and run a series of tests against them.

#### Prerequisite

Create a Gitlab personnal access token with the following scopes:
- `api`
- `read_api`
- `write_repository`

Create a GitHub App with the following permissions:
- `administration read/write`
- `checks read/write`
- `commit statuses read/write`
- `contents read/write`
- `issues read/write`
- `metadata read-only`

#### Run integration tests

```bash
$ export GITHUB_APP_ID=<github app id>
$ export GITHUB_AUTHOR=<github author, usually username[bot]>
$ export GITHUB_INSTALLATION_ID=<github installation id>
$ export GITHUB_OWNER=<github repository owner>
$ export GITHUB_PRIVATE_KEY_PATH=<github private key path>
$ export GITLAB_AUTHOR=<gitlab author, usually owner>
$ export GITLAB_OWNER=<gitlab repository owner>
$ export GITLAB_TOKEN=<gitlab api token>
$ make integration
```

## Release

[GoReleaser](https://goreleaser.com/) is used to generate all the necessary
binaries and attach them together with the changelog to the GitHub release. To
release, create a tag then wait for GitHub Actions to publish the release.
