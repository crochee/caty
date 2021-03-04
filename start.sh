#!/bin/bash
set -ex

export config_path=conf/config.yml

go build -tags=jsoniter
chmod +x ./obs
./obs >>console.txt 2 & 1 &