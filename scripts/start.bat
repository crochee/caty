@echo off
go build -ldflags="-s -w" -tags jsoniter -o obs.exe ./cmd/server
obs.exe