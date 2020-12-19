// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package bucket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	BucketName string `json:"bucket_name" binding:"required"`
}

func CreateBucket(ctx *gin.Context) {
	var createRequest CreateRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, 1)
		return
	}
}
