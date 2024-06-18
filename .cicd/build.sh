#!/bin/sh
set -eu

set -a
. .cicd/env
. .cicd/functions.sh
set +a

echo Building $APP_NAME-$APP_COMPONENT

export GOPATH='/woodpecker/go'
export CGO_ENABLED=0

cd "./$APP_COMPONENT"
go build -v -a -o ../bin/app
ls -lh ../bin

echo Done
