// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package file

import "mime/multipart"

type Name struct {
	BucketName string `json:"bucket_name" uri:"bucket_name" form:"bucket_name" binding:"required"`
}

type Id struct {
	BucketId uint `json:"bucket_id" uri:"bucket_id" form:"bucket_id" binding:"required"`
}

type Info struct {
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type Target struct {
	Id
	FileId uint `json:"file_id" uri:"file_id" form:"file_id" binding:"required"`
}
