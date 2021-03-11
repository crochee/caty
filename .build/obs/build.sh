#!/bin/bash
set -ex

WORKER_SPACE=$(pwd)
COMPONENT=osb

function compile() {
  go build -o ${COMPONENT} -tags jsoniter .
  cp ${COMPONENT} ${WORKER_SPACE}/.build/obs
  cp -r ${WORKER_SPACE}/conf ${WORKER_SPACE}/.build/obs
  mkdir -p ${WORKER_SPACE}/.build/obs/docs
  cp -r ${WORKER_SPACE}/docs/swagger* ${WORKER_SPACE}/.build/obs/docs/
}

function images() {
  docker build -t ${COMPONENT}:latest .
  docker tag ${COMPONENT}:latest obs:latest
}

function clean() {
  rm -rf ${WORKER_SPACE}/.build/obs/conf
  rm -rf ${WORKER_SPACE}/.build/obs/${COMPONENT}
  rm -rf ${WORKER_SPACE}/.build/obs/docs
}

compile
images
clean
