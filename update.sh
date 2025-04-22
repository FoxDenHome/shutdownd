#!/bin/bash
set -euo pipefail
set -x

# This needs to run as root probably

cd "$(realpath $(dirname "${0}"))"
git pull
./build.sh

useradd -s /bin/false shutdownd || true

cp -fv shutdownd.service /etc/systemd/system/shutdownd.service
cp -fv shutdownd.sudoers /etc/sudoers.d/shutdownd
systemctl daemon-reload
systemctl enable shutdownd
systemctl restart shutdownd
