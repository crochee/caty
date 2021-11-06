#!/bin/sh

set -ex

# docker build --no-cache -t caty .
# docker run -itd -p 8120:8120 --restart=always --name catysrv caty
if [ "${1:0:1}" = '-' ]; then
  set -- caty "$@"
fi

exec "$@"
