#!/bin/bash
set -euo pipefail

export COMMIT=$(git rev-parse --short HEAD)

buildos() {
    export GOOS="$1"
    export GOARCH="$2"
    export CGO_ENABLED=0

    local suffix=''
    if [[ "${GOOS}" == "windows" ]]; then
        suffix='.exe'
    fi

    go build -o "dist/shutdownd-${GOOS}-${GOARCH}${suffix}" \
        -ldflags "-s -w -X util.commit=${COMMIT}" \
        -trimpath \
        -v \
        "./cmd/shutdownd"
}

rm -rf dist

go mod download

buildos linux amd64
buildos linux arm64
buildos darwin arm64
buildos windows amd64
