// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/31

package main

import "github.com/asim/go-micro/v3/api/handler/web"

func ser() {
	web.NewHandler()
	web.WithService()
}
