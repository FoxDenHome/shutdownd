#!/bin/bash
set -euo pipefail
set -x

id shutdownd >/dev/null || useradd -r -d /var/empty -s /usr/sbin/nologin shutdownd

_path="$(realpath $(dirname "${0}"))"
cd "$_path"

cp -fv shutdownd /usr/bin/shutdownd

mkdir -p /etc/shutdownd && chown -R shutdownd:shutdownd /etc/shutdownd

cp -fv shutdownd.service /etc/systemd/system/shutdownd.service
cp -fv shutdownd.sudoers /etc/sudoers.d/shutdownd
systemctl daemon-reload
systemctl enable shutdownd
systemctl restart shutdownd
