#!/bin/bash
set -e

export OS_DIR_CONFIG="/etc/mysterium-node"
export OS_DIR_DATA="/var/lib/mysterium-node"
export OS_DIR_RUN="/var/run/mysterium-node"

/usr/local/bin/prepare-run-env.sh

exec /usr/bin/myst \
 --config-dir=$OS_DIR_CONFIG \
 --script-dir=$OS_DIR_CONFIG \
 --log-dir="" \
 --data-dir=$OS_DIR_DATA \
 --runtime-dir=$OS_DIR_RUN \
 $@
