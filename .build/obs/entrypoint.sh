#!/bin/bash
set -ex

if [ "${1:0:1}" = '-' ]; then
  set -- obs "$@"
fi

exec "$@"
