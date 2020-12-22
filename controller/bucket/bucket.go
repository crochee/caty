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
	"obs/util"
)

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param request body model.CreateBucket true "bucket"
// @Success 200 {object} model.AkSk
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket [post]
func CreateBucket(ctx *gin.Context) {
	var createBucket model.CreateBucket
	if err := ctx.ShouldBindJSON(&createBucket); err != nil {
		logger.Errorf("bind request failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
		return
	}
	token := verify.NewToken(createBucket.BucketName.BucketName)
	token.AddAction(createBucket.Action)
	ak, sk, err := token.Create()
	if err != nil {
		logger.Errorf("create token failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "create bucket failed")
		return
	}
	path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + createBucket.BucketName.BucketName
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
// @Param request body model.BucketName true "bucket"
// @Success 200
// @Header 200 {string} Last-Modified "last modify"
// @Failure 404
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket [head]
func HeadBucket(ctx *gin.Context) {
	var bucket model.BucketName
	if err := ctx.ShouldBindJSON(&bucket); err != nil {
		logger.Errorf("bind request failed.Error:%v", err)
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
		logger.Errorf("get %s stat failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "find bucket failed")
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
// @Param request body model.SimpleBucket true "bucket"
// @Success 200
// @Failure 404
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket [delete]
func DeleteBucket(ctx *gin.Context) {
	var simpleBucket model.SimpleBucket
	if err := ctx.ShouldBindJSON(&simpleBucket); err != nil {
		logger.Errorf("bind request failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
		return
	}
	t := verify.Token{
		AkSecret: util.Slice(simpleBucket.Ak),
		Bucket:   simpleBucket.BucketName.BucketName,
	}
	if err := t.Verify(simpleBucket.Sk); err != nil {
		logger.Errorf("verify sk failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "check your ak sk"))
		return
	}
	for action := range t.Action {
		if action >= verify.Delete {
			path := config.Cfg.YamlConfig.ServiceInformation.SaveRootPath + simpleBucket.BucketName.BucketName
			if err := os.RemoveAll(path); err != nil {
				logger.Errorf("delete path(%s) failed.Error:%v", path, err)
				response.ErrorWithMessage(ctx, "delete bucket failed")
				return
			}
			ctx.Status(http.StatusOK)
			return
		}
	}
	ctx.Status(http.StatusNotFound)
}
