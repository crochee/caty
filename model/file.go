// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/24

package model

import "mime/multipart"

type FileInfo struct {
	FileTarget
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type FileTarget struct {
	Path string `json:"path" uri:"path" form:"path" binding:"required"`
}
