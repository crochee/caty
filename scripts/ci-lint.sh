#!/bin/bash
set -e

golangci-lint run -c build/ci/golangci-lint-demo.yml --tests=false  --out-format=json  > golangci-lint.json 2>&1