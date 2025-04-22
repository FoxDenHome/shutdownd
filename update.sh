#!/bin/bash
set -euo pipefail
set -x

cd /opt/shutdownd
git pull
./build.sh

sudo cp -fv shutdownd.service /etc/systemd/system/shutdownd.service
sudo cp -fv shutdownd.sudoers /etc/sudoers.d/shutdownd
sudo systemctl daemon-reload
sudo systemctl enable shutdownd
sudo systemctl restart shutdownd
