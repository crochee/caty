// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/22

package model

import "obs/service/verify"

type BucketAction struct {
	Action verify.BucketAction `json:"action"`
}

type BucketName struct {
	BucketName string `json:"bucket_name" uri:"bucket_name" form:"bucket_name" binding:"required"`
}
