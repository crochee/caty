@echo off
go build -ldflags="-s -w" -tags jsoniter -o obs.exe ./cmd/server
set config=conf/config.yml
obs.exe