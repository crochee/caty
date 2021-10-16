#!/bin/bash
set -ex
# docker build --no-cache -t cca .
# docker run -itd -p 8120:8120 --restart=always --name cca-server bash
if [ "${1:0:1}" = '-' ]; then
  set -- cca "$@"
fi

exec "$@"
