#!/bin/sh
set -eu

set -a
. .ci/lib.sh
set +a

echo && echo "Testing $APP_NAME-$APP_COMPONENT"

export GOPROXY="$GO_PROXY,https://proxy.golang.org,direct"
export GOPATH='/woodpecker/go'

cd "./$APP_COMPONENT"
go test -cover -v ./...

echo && echo 'Done'
