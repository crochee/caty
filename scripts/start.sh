#!/bin/bash
set -ex

export config_path=conf/config.yml
# 安全编译选项
export CGO_ENABLED=1
export CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -Wformat=2 -Wfloat-equal -Wshadow"
export CGO_LDFLAGS="-Wl,-z,relro,-z,now"
# go build
go build -trimpath -ldflags="-s -w" -buildmode=pie -o=obs -tags=jsoniter ./cmd/server

chmod +x ./obs
./obs >>console.txt 2>&1 &
