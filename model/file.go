// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/24

package model

import "mime/multipart"

type FileInfo struct {
	SimpleBucket
	Target string                `json:"target" form:"target" binding:"required"`
	File   *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}
