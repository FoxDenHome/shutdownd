#!/bin/bash
set -euo pipefail

export COMMIT="$(git rev-parse HEAD)"

buildbin() {
    go build \
        -o "dist/" \
        -ldflags "-s -w -X github.com/FoxDenHome/shutdownd/util.commit=${COMMIT}" \
        -trimpath \
        "./cmd/$1"
}

rm -rf dist && mkdir -p dist

buildbin shutdownd
buildbin certgen
