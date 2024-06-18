#!/bin/sh
set -eu

set -a
. .cicd/env
. .cicd/functions.sh
set +a

echo Testing $APP_NAME-$APP_COMPONENT

export GOPATH='/woodpecker/go'

cd "./$APP_COMPONENT"
go test -cover -v ./...

echo Done
