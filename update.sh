#!/bin/bash
set -euo pipefail
set -x

# This needs to run as root probably

_path="$(realpath $(dirname "${0}"))"
cd "$_path"
git pull
./build.sh

useradd -s /bin/false shutdownd || true

sed "s~__PATH__~$_path~g" shutdownd.service /etc/systemd/system/shutdownd.service
cp -fv shutdownd.sudoers /etc/sudoers.d/shutdownd
systemctl daemon-reload
systemctl enable shutdownd
systemctl restart shutdownd
