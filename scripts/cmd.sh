#!/bin/bash

set -ex

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build ./cmd/catyctl