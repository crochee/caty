// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

type Name struct {
	BucketName string `json:"bucket_name" uri:"bucket_name" form:"bucket_name" binding:"required"`
}

type Id struct {
	BucketId uint `json:"bucket_id" uri:"bucket_id" form:"bucket_id" binding:"required"`
}
