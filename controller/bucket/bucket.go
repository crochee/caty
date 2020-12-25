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
	"obs/model"
	"obs/response"
	"obs/service/verify"
)

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param request body model.BucketAction true "bucket action"
// @Success 200 {object} model.AkSk
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name} [post]
func CreateBucket(ctx *gin.Context) {
	var bucket model.BucketName
	if err := ctx.ShouldBindUri(&bucket); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "bucket_name is nil"))
		return
	}
	var bucketAction model.BucketAction
	if err := ctx.ShouldBindJSON(&bucketAction); err != nil {
		logger.Errorf("bind body failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
		return
	}
	newToken := &verify.Token{}
	newToken.AddAction(bucketAction.Action)
	ak, sk, err := newToken.Create()
	if err != nil {
		logger.Errorf("create token failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "create bucket failed")
		return
	}
	path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + bucket.BucketName
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Errorf("mkdir %s failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "create bucket failed")
		return
	}
	logger.Info("create bucket success")
	ctx.JSON(http.StatusOK, model.AkSk{
		Ak: ak,
		Sk: sk,
	})
}

// HeadBucket godoc
// @Summary HeadBucket
// @Description head bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Success 200
// @Header 200 {string} Last-Modified "last modify"
// @Failure 404
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name} [head]
func HeadBucket(ctx *gin.Context) {
	var bucket model.BucketName
	if err := ctx.ShouldBindUri(&bucket); err != nil {
		logger.Errorf("bind uri failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "bucket_name is nil"))
		return
	}
	path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + bucket.BucketName
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.Status(http.StatusNotFound)
			return
		}
		logger.Errorf("find path(%s) failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "get bucket info failed")
		return
	}
	if !fileInfo.IsDir() {
		ctx.Status(http.StatusNotFound)
		return
	}
	logger.Infof("get bucket(%s) status success", path)
	ctx.Header("Last-Modified", fileInfo.ModTime().String())
	ctx.Status(http.StatusOK)
}

// DeleteBucket godoc
// @Summary DeleteBucket
// @Description delete bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Success 200
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name} [delete]
func DeleteBucket(ctx *gin.Context) {
	var bucket model.BucketName
	if err := ctx.ShouldBindUri(&bucket); err != nil {
		logger.Errorf("bind uri  failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "bucket_name is nil"))
		return
	}
	path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + bucket.BucketName
	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("delete path(%s) failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "delete bucket failed")
		return
	}
	ctx.Status(http.StatusOK)
}
