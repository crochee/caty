#!/bin/bash

set -ex

image=docker.io/evildecay/etcdkeeper:latest
if [ -n "$(docker images -q ${image})" ]; then
  docker pull ${image}
fi

docker run -itd -p 8080:8080 --restart=always --name etcdkeeper ${image}
