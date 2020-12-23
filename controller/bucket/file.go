// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package bucket

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(ctx *gin.Context) {
	http.ServeFile()
}
