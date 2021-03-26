@echo off
golangci-lint run .\... -c golangci-lint.yml --tests=false  --out-format=json  > golangci-lint.json 2>&1