#!/bin/bash
set -euo pipefail
set -x

cd "$(realpath $(dirname "${0}"))"
git pull
git checkout main
git reset --hard origin/main
git clean -fdx
./build.sh
sudo ./install.sh
