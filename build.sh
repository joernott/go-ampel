#!/bin/bash
set -e
GO=$(command -v go)
DOCKER=$(command -v docker)
TAG=$1

if [ -z "${GO}" ]; then
    echo "You need to install golang to be able to build this."
    exit 10
fi
if [ -z "${GOPATH}" ]; then
    echo "GOPATH not set"
fi
if [ -z "${DOCKER}" ]; then
    echo "Docker is not available, not building a docker image."
fi

$GO get -v
CGO_ENABLED=0 GOOS=linux $GO build -v -installsuffix cgo -o ampel .

if [ -n "${DOCKER}" ]; then
    if [ -z "${TAG}" ]; then
        TAG="ampel"
    fi
    $DOCKER build -t ${TAG} .
fi
