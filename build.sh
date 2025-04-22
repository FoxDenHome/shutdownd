#!/bin/bash
set -euo pipefail

export COMMIT="$(git rev-parse HEAD)"

buildbin() {
    go build \
        -ldflags "-s -w -X github.com/FoxDenHome/shutdownd/util.commit=${COMMIT}" \
        -trimpath \
        "./cmd/$1"
}

rm -fv shutdownd certgen

buildbin shutdownd
buildbin certgen
