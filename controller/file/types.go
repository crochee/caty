// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package file

import "mime/multipart"

type Name struct {
	BucketName string `json:"bucket_name" uri:"bucket_name" form:"bucket_name" binding:"required"`
}

type Info struct {
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type Target struct {
	Name
	FileName string `json:"file_name" uri:"file_name" form:"file_name" binding:"required"`
}
