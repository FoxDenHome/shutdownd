#!/bin/bash
set -euo pipefail

export COMMIT=$(git rev-parse --short HEAD)

buildbin() {
    local suffix=''
    if [[ "${GOOS}" == "windows" ]]; then
        suffix='.exe'
    fi

    local prefix=''
    if [[ "$1" != "shutdownd" ]]; then
        prefix='shutdownd-'
    fi

    go build -o "dist/$prefix$1-${GOOS}-${GOARCH}${suffix}" \
        -ldflags "-s -w -X util.commit=${COMMIT}" \
        -trimpath \
        -v \
        "./cmd/$1"
}

buildos() {
    export GOOS="$1"
    export GOARCH="$2"
    export CGO_ENABLED=0

    buildbin certgen
    buildbin shutdownd
}

rm -rf dist

go mod download

buildos linux amd64
buildos linux arm64
buildos darwin arm64
buildos windows amd64
