// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package bucket

import (
	"github.com/gin-gonic/gin"

	"obs/response"
)

type CreateRequest struct {
	BucketName string `json:"bucket_name" binding:"required"`
}

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param request body CreateRequest true "bucket"
// @Success 201
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file [post]
func CreateBucket(ctx *gin.Context) {
	var createRequest CreateRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		response.ErrorWith(ctx, err)
		return
	}

}
