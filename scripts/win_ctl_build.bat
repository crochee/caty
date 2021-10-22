@echo off
set GOOS=linux
set CGO_ENABLED=0
set GOARCH=amd64
go build ./cmd/ccactl