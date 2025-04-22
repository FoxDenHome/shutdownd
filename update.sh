#!/bin/bash
set -euo pipefail
set -x

_path="$(realpath $(dirname "${0}"))"
cd "$_path"
git pull
./build.sh
./install.sh
