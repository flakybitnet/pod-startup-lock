#!/bin/sh
set -eu

set -a
. .cicd/env
. .cicd/functions.sh
set +a

DST_REGISTRY='public.ecr.aws'
DST_REPOSITORY="$APP_NAME/$APP_COMPONENT"
ECR_CREDS=$(printf '%s:%s' "AWS" "$(cat $AWS_PWD_FILE)")

SRC_IMAGE="$HARBOR_REGISTRY/$HARBOR_PROJECT/$HARBOR_REPOSITORY:$APP_VERSION"
DST_IMAGE="$DST_REGISTRY/$NAMESPACE/$DST_REPOSITORY:$APP_VERSION"

echo Publishing $DST_IMAGE image
retry 2 skopeo copy --dest-creds="$ECR_CREDS" "docker://$SRC_IMAGE" "docker://$DST_IMAGE"

SRC_IMAGE="$HARBOR_REGISTRY/$HARBOR_PROJECT/$HARBOR_REPOSITORY:$APP_VERSION-debug"
DST_IMAGE="$DST_REGISTRY/$NAMESPACE/$DST_REPOSITORY:$APP_VERSION-debug"

echo Publishing $DST_IMAGE image
retry 2 skopeo copy --dest-creds="$ECR_CREDS" "docker://$SRC_IMAGE" "docker://$DST_IMAGE"

echo Done
