#!/bin/bash
set -ex

export config=./conf/config.yml
if [ "${1:0:1}" = '-' ]; then
  set -- obs "$@"
fi

exec "$@"
