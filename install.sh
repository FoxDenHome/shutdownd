#!/bin/bash
set -euo pipefail
set -x

id shutdownd >/dev/null || useradd -r -d /var/empty -s /usr/sbin/nologin shutdownd

_path="$(realpath $(dirname "${0}"))"
cd "$_path"

go build -trimpath ./cmd/shutdownd
cp -fv shutdownd /usr/bin/shutdownd

mkdir -p /etc/shutdownd && chown -R shutdownd:shutdownd /etc/shutdownd

cp -fv *.service /etc/systemd/system/
cp -fv 10-shutdownd.rules /etc/polkit-1/rules.d/
systemctl daemon-reload
systemctl restart polkit
systemctl enable shutdownd
systemctl restart shutdownd
