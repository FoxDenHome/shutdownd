#!/bin/bash
set -euo pipefail
set -x

# This needs to run as root
useradd -s /bin/false shutdownd || true

_path="$(realpath $(dirname "${0}"))"
cd "$_path"
git pull
./build.sh

cat shutdownd.service | sed "s~__PATH__~$_path~g" > /etc/systemd/system/shutdownd.service
cp -fv shutdownd.sudoers /etc/sudoers.d/shutdownd
systemctl daemon-reload
systemctl enable shutdownd
systemctl restart shutdownd
