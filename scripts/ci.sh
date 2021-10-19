#!/bin/bash

set -ex

if [[ -n $(docker ps -a |grep ccasrv) ]]; then
  docker rm ccasrv --force
fi

if [[ -n $(docker images -q cca:latest) ]]; then
  docker rmi cca:latest
fi

docker build -f ./build/cca/Dockerfile -t cca:latest .
docker run -itd -p 8120:8120 --restart=always --name ccasrv cca:latest
