#!/bin/bash

set -ex

if [[ -n $(docker ps -a |grep catysrv) ]]; then
  docker rm catysrv --force
fi

if [[ -n $(docker images -q caty:latest) ]]; then
  docker rmi caty:latest
fi

docker build --no-cache -f ./build/caty/Dockerfile -t caty:latest .
docker run -itd -p 8120:8120 --restart=always --name catysrv caty:latest