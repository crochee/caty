#!/bin/bash

set -ex

GO_VERSION=$(go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
if [[ "$GO_VERSION" -eq "" ]]; then
  echo "please install golang"
  exit 1
fi

VERSION1380="v1.38.0"
VERSION1421="v1.42.1"

if [[ $GO_VERSION -gt 16 ]]; then
  hash golangci-lint >/dev/null 2>&1
  OP_MODE=$?
  if [[ $OP_MODE -ne 0 ]]; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1 >/dev/null 2>&1
  fi
  CI_VERSION=$(golangci-lint --version)
  if [[ $CI_VERSION != *$VERSION1421* ]]; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1 >/dev/null 2>&1
  fi
  golangci-lint run -c ./build/lint/golangci-lint-demo.yml --tests=false --out-format=json >golangci-lint.json 2>&1
else
  hash golangci-lint >/dev/null 2>&1
  OP_MODE=$?
  if [[ $OP_MODE -ne 0 ]]; then
    go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0 >/dev/null 2>&1
  fi
  CI_VERSION=$(golangci-lint --version)
  if [[ $CI_VERSION != *$VERSION1380* ]]; then
    go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0 >/dev/null 2>&1
  fi
  golangci-lint run -c ./build/lint/golangci-lint.yml --tests=false --out-format=json >golangci-lint.json 2>&1
fi
