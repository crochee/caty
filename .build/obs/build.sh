#!/bin/bash
set -ex

WORKER_SPACE=$(pwd)
COMPONENT=osb

component_compile() {
  go build -o ${COMPONENT} -tags jsoniter .
  cp ${COMPONENT} ${WORKER_SPACE}/.build/obs
  cp -r ${WORKER_SPACE}/conf ${WORKER_SPACE}/.build/obs
  mkdir -p ${WORKER_SPACE}/.build/obs/docs
  cp -r ${WORKER_SPACE}/docs/swagger* ${WORKER_SPACE}/.build/obs/docs/
}

docker_build() {
  docker build -t ${COMPONENT}:latest .
  docker tag ${COMPONENT}:latest obs:latest
}

clean_file() {
  rm -rf ${WORKER_SPACE}/.build/obs/conf
  rm -rf ${WORKER_SPACE}/.build/obs/${COMPONENT}
  rm -rf ${WORKER_SPACE}/.build/obs/docs
}

component_compile
docker_build
clean_file
