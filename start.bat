@echo off
go build -tags jsoniter -o obs.exe
set config=conf/config.yml
obs.exe