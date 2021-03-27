// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/25

package main

//go:generate go install github.com/swaggo/swag/cmd/swag@v1.7.0
//go:generate swag i -g router/router.go

//go:generate go install github.com/securego/gosec/v2/cmd/gosec@v2.7.0
//go:generate gosec -fmt=json -out=results.json .\...

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0
//go:generate golangci-lint run .\... -c golangci-lint.yml --tests=false –-out-format=json > golangci-lint.json 2>&1