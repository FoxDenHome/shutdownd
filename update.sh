#!/bin/bash
set -euo pipefail
set -x

cd "$(realpath $(dirname "${0}"))"
git pull
./build.sh
sudo ./install.sh
