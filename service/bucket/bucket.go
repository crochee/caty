// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"

	"obs/config"
	"obs/logger"
	"obs/model/db"
	"obs/response"
	"obs/service/tokenx"
	"obs/util"
)

// CreateBucket 创建桶
//
// @param ctx 请求context
// @param token token信息
// @param bucketName 桶名
// @Success uint 桶id
// @Failure error 自定义错误
func CreateBucket(ctx context.Context, token *tokenx.Token, bucketName string) (uint, error) {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	bucket := &db.Bucket{
		Bucket: bucketName,
		Domain: token.Domain,
		User:   token.User,
	}
	if err := tx.Create(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("insert db failed.Error:%v", err)
		return 0, response.Errors(http.StatusInternalServerError, err)
	}
	path, err := filepath.Abs(fmt.Sprintf("%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return 0, response.Error(http.StatusInternalServerError, "clear bucket failed")
	}
	if err = os.MkdirAll(filepath.Clean(path), os.ModePerm); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("create bucket failed.Error:%v", err)
		return 0, response.Error(http.StatusInternalServerError, "create bucket failed")
	}
	return bucket.ID, nil
}

// HeadBucket 查询桶信息
//
// @param ctx 请求context
// @param token token信息
// @param bucketId 桶名
// @Success *Info 桶信息
// @Failure error 自定义错误
func HeadBucket(ctx context.Context, token *tokenx.Token, bucketId uint) (*Info, error) {
	conn := db.NewDB()
	bucket := &db.Bucket{ID: bucketId}
	if err := conn.Model(bucket).Where("id =? AND domain= ?", bucketId, token.Domain).Find(bucket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return nil, response.Errors(http.StatusInternalServerError, err)
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logger.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return nil, response.Error(http.StatusInternalServerError, "get bucket failed")
	}
	path = filepath.Clean(path)

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			logger.FromContext(ctx).Errorf("not find bucket %s failed.Error:%v", bucket.Bucket, err)
			return nil, response.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
		}
		logger.FromContext(ctx).Errorf("find path(%s) failed.Error:%v", path, err)
		return nil, response.Error(http.StatusInternalServerError, "get bucket failed")
	}
	var info = &Info{
		LastModified: fileInfo.ModTime(),
		Name:         fileInfo.Name(),
	}
	if info.Size, info.Count, err = util.DirSize(path); err != nil {
		logger.FromContext(ctx).Errorf("lookup path %s failed.Error:%v", path, err)
		return nil, response.Error(http.StatusInternalServerError, "find bucket's all file failed")
	}
	return info, nil
}

// DeleteBucket 删除桶
//
// @param ctx 请求context
// @param token token信息
// @param bucketId 桶名
// @Failure error 自定义错误
func DeleteBucket(ctx context.Context, token *tokenx.Token, bucketId uint) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()

	bucket := &db.Bucket{ID: bucketId}
	if err := tx.Model(bucket).Where("id =? AND domain= ?", bucketId, token.Domain).Find(bucket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			return response.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
		}
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	if err := tx.Delete(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("delete bucket failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return response.Error(http.StatusInternalServerError, "get bucket failed")
	}
	path = filepath.Clean(path)
	if err = os.RemoveAll(path); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("remove %s failed.Error:%v", path, err)
		return response.Error(http.StatusInternalServerError, "delete bucket failed")
	}
	return nil
}
