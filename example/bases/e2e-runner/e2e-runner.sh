#!/usr/bin/env bash

set -e

# mandatory variables
: ${E2E_FLOW?"Environment variable E2E_FLOW is undefined"}
: ${KUBERNETES_NAMESPACE?"Environment variable KUBERNETES_NAMESPACE is undefined"}

# wait for deployments to be created and available
kubectl wait --namespace="$KUBERNETES_NAMESPACE" \
  --for=condition=Available \
  --selector "app in (backend, frontend)" \
  --timeout=300s \
  deployments

# wait for readiness
echo "Waiting for frontend pods to be ready..."
kubectl wait --namespace="$KUBERNETES_NAMESPACE" \
  --for=condition=Ready \
  --selector "app in (backend, frontend)" \
  --timeout=300s \
  pods

# run e2e, currently only a placeholder for the tests
echo "e2e tests execution was successful"

# list all deployments
image_list=$(kubectl get deploy -o jsonpath='{.items[*].spec.template.spec.containers[*].image}')

# set commit status in the corresponding repository/commit sha
for image in $image_list
do
  echo "Getting metadata for ${image}"
  repository=$(docker inspect $image --format='{{.Config.Labels.repository}}')
  commitSha=$(docker inspect $image --format='{{.Config.Labels.commitSha}}')
  if [[ "$repository" == "<no value>" ]] || [[ "$commitSha" == "<no value>" ]]
  then
    echo "Label repository or label commit sha are undefined. Skipping..."
    continue
  fi
  echo "Found ${repository} with sha ${commitSha}"

  grgate status set "$repository" \
    --commit "$commitSha" \
    --name "$E2E_FLOW" \
    --status completed \
    --state success
done
