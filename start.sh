#!/bin/bash
set -ex

export config_path=conf/config.yml

go build -ldflags="-s -w" -tags jsoniter -o obs cmd/server
chmod +x ./obs
./obs >>console.txt 2 & 1 &