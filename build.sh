#!/bin/bash
set -euo pipefail
set -x

export COMMIT="$(git rev-parse HEAD)"

go build -o . \
    -ldflags "-s -w -X github.com/FoxDenHome/shutdownd/util.commit=${COMMIT}" \
    -trimpath \
    "./cmd/..."
