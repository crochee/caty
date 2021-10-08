#!/bin/bash
set -ex
# docker build --no-cache -t cca .
# docker run -itd -p 8150:8150 -v /home/lcf/cloud/data:/cca/ --restart=always --name ccav1 cca
if [ "${1:0:1}" = '-' ]; then
  set -- cca "$@"
fi

exec "$@"
