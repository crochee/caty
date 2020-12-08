// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/8

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"obs/logger"
)

type Bucket struct {
	Name string `json:"name"`
}

// BuildBucket godoc
// @Summary BuildBucket
// @Description build bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param request body Bucket true "bucket content"
// @Success 201
// @Failure 400
// @Failure 500 {object} util.HttpErr
// @Router /user/login [post]
func BuildBucket(ctx *gin.Context) {
	var bucket Bucket
	if err := ctx.ShouldBindBodyWith(&bucket, binding.JSON); err != nil {
		logger.Errorf("build bucket bind body failed.Error:%v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

}
