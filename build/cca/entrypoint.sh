#!/bin/bash
set -ex
# docker build --no-cache -t cca .
# docker run -itd -p 8120:8120 --restart=always --name ccasrv cca
if [ "${1:0:1}" = '-' ]; then
  set -- cca "$@"
fi

exec "$@"
