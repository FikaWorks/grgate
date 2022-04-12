GRGate example with Kubernetes and GitHub
=========================================

> This repository contains an example on how to run e2e tests for each
> deployment and automatically publish the corresponding releases when the test
> execution is successful. Docker KinD is used, and you will need to create 2
> repositories to test the release publishing flow.

## Overview

The diagram below represent the idea of what GRGate can be used in Kubernetes.

![GRGate in your GitOps environment](grgate-kubernetes.png "GRGate in your
GitOps environment")

In the following section, we will test the integration of GRGate with Docker
KinD.

## Getting started

1. locally run [Kubernetes in
   Docker](https://docs.docker.com/desktop/kubernetes/) (KinD)
1. **prepare repositories**: create 1 or more repositories (can be private), in
   our example we will use: frontend and backend. In both repositories, add the
   following Dockerfile:

   ```Dockerfile
   FROM nginx:alpine
   ```
   Note: this is an example, usually you would have an application doing stuff
   there

   Then add the following `.grgate.yaml` file:

   ```yaml
   enabled: true
   statuses:
     - e2e happy flow
   ```

1. **build containers**: from those repositories locally build both frontend
   and backend containers with the corresponding labels. For example, the
   Docker build command for the frontend can be:

   ```bash
   docker build \
     --label repository="my-org/frontend" \
     --label commitSha="$(git rev-parse HEAD)" \
     -t frontend \
     .
   ```

   And our backend:

   ```bash
   docker build \
     --label repository="my-org/backend" \
     --label commitSha="$(git rev-parse HEAD)" \
     -t backend \
     .
   ```
1. **install GRGate app** from the
   [Marketplace](https://github.com/marketplace/grgate/):
1. **create draft releases**: for each repositories you previously created,
   use the gh cli to create draft releases:

   ```bash
   gh release create v1.0.0 --draft --generate-notes
   ```
1. **prepare manifests**: create a new GitHub App with the following
   permissions:
   - `checks read/write`
   - `commit statuses read/write`
   - `metadata read-only`

   The app is used by the e2e runner to set the status check to the target
   commit defined in the Docker image.

   Copy the GitHub App ID and installation ID to the
   `overlays/development/grgate-config.yaml`.

   Generate a certificate and store it in
   `overlays/development/github.private-key.pem`.

   Finally, install the app in your organization.

1. **deploy**:

   The following make command build the e2e-runner Docker image and deploy it
   to KinD:
   ```bash
   make build deploy-development
   ```

1. **validate**: after the e2e tests are executed, commit status should be
   updated in both repositories and GRGate should have published the releases.
