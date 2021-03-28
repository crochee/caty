#!/bin/bash
set -ex
# docker run -itd -p 8150:8150 -v /home/lcf/cloud/data:/obs/ --restart=always --name obsv1 obs
if [ "${1:0:1}" = '-' ]; then
  set -- obs "$@"
fi

exec "$@"
