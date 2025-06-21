#!/bin/bash

set -e

source .env
[ -n "$VERSION" ] || VERSION='1.0.0'
[ -n "$DATE" ] || DATE=$(date +%Y-%m-%d)
export VERSION DATE
PKGS=$(go list ./... | grep -v /vendor/)

go mod tidy
go generate $PKGS
go vet $PKGS
go fmt $PKGS
go test -race $PKGS

cd cmd/server
go build -ldflags "\
    -X 'main.buildVersion=${VERSION}' \
    -X 'main.buildDate=${DATE}' \
    " .
cd ../..

# cd cmd/client
# go build -ldflags "\
#     -X 'main.buildVersion=${VERSION}' \
#     -X 'main.buildDate=${DATE}' \
#     ".
# cd ../..

docker compose build goph-keeper
