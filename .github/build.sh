#!/bin/bash
set -euo pipefail

export COMMIT="$(git rev-parse HEAD)"
export CGO_ENABLED=0

buildbin() {
    local bin="$1"

    local suffix=''
    if [[ "${GOOS}" == "windows" ]]; then
        suffix='.exe'
    fi

    local prefix=''
    if [[ "${bin}" != "shutdownd" ]]; then
        prefix='shutdownd-'
    fi

    go build \
        -o "dist/${prefix}${bin}-${GOOS}-${GOARCH}${suffix}" \
        -ldflags "-s -w -X github.com/FoxDenHome/shutdownd/util.commit=${COMMIT}" \
        -trimpath \
        "./cmd/${bin}"
}

buildos() {
    export GOOS="$1"
    export GOARCH="$2"

    buildbin certgen
    buildbin shutdownd
}

rm -rf dist

go mod download

buildos linux amd64
buildos linux arm64
buildos darwin arm64
buildos windows amd64
