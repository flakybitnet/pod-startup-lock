#!/bin/sh
set -eu

set -a
. .ci/lib.sh
set +a

echo && echo "Building $APP_NAME-$APP_COMPONENT"

export GOPROXY="$GO_PROXY,https://proxy.golang.org,direct"
export GOPATH='/woodpecker/go'
export CGO_ENABLED=0

cd "./$APP_COMPONENT"
go build -v -a -o ../bin/app

echo 'bin:'
ls -lh ../bin

echo && echo 'Done'
