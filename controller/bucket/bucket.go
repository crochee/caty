// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package bucket

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"obs/config"
	"obs/logger"
	"obs/response"
	"obs/service/verify"
)

type CreateRequest struct {
	BucketName string              `json:"bucket_name" binding:"required"`
	Action     verify.BucketAction `json:"action"`
}

type BucketAkSk struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param request body CreateRequest true "bucket"
// @Success 200 {object} BucketAkSk
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket [post]
func CreateBucket(ctx *gin.Context) {
	var createRequest CreateRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	token := verify.NewToken(createRequest.BucketName)
	token.AddAction(createRequest.Action)
	ak, sk, err := token.Create()
	if err != nil {
		logger.Errorf("create token failed.Error:%v", err)
		response.ErrorWith(ctx, err)
		return
	}
	path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + createRequest.BucketName
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Errorf("mkdir %s failed.Error:%v", path, err)
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, BucketAkSk{
		Ak: ak,
		Sk: sk,
	})
}
