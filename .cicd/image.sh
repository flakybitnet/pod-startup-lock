#!/bin/sh
set -eu

set -a
. .cicd/env
. .cicd/functions.sh
set +a

IMAGE="$HARBOR_PROJECT/$HARBOR_REPOSITORY:$APP_VERSION"
DOCKERFILE=".docker/Dockerfile"

if [ "${IMAGE_DEBUG:-false}" = "true" ]; then
  echo Debug image is set to buid
  IMAGE="$IMAGE-debug"
  DOCKERFILE="$DOCKERFILE-debug"
fi

echo Building $IMAGE image

executor --context ./ \
    --dockerfile "$DOCKERFILE" \
    --destination "$HARBOR_REGISTRY/$IMAGE"

echo Done

