// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/22

package model

import "obs/service/verify"

type CreateBucket struct {
	BucketName
	Action verify.BucketAction `json:"action"`
}

type AkSk struct {
	Ak string `json:"ak" binding:"required"`
	Sk string `json:"sk" binding:"required"`
}

type SimpleBucket struct {
	BucketName
	AkSk
}

type BucketName struct {
	BucketName string `json:"bucket_name" binding:"required"`
}
