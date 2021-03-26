// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"obs/config"
	"obs/e"
	"obs/logger"
	"obs/model/db"
	"obs/service/tokenx"
	"obs/util"
)

// CreateBucket 创建桶
func CreateBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	bucket := &db.Bucket{
		Bucket: bucketName,
		Domain: token.Domain,
	}
	if err := tx.Create(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("insert db failed.Error:%v", err)
		return e.New(e.Unknown, err.Error())
	}
	path, err := filepath.Abs(fmt.Sprintf("%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return e.Error(http.StatusInternalServerError, "clear bucket failed")
	}
	if err = os.MkdirAll(filepath.Clean(path), os.ModePerm); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("create bucket failed.Error:%v", err)
		return e.Error(http.StatusInternalServerError, "create bucket failed")
	}
	return nil
}

// HeadBucket 查询桶信息 Info
func HeadBucket(ctx context.Context, token *tokenx.Token, bucketName string) (*Info, error) {
	conn := db.NewDB()
	bucket := &db.Bucket{Bucket: bucketName}
	if err := conn.Find(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return nil, response.Errors(http.StatusInternalServerError, err)
	}
	if bucket.Domain != token.Domain {
		return nil, e.Error(http.StatusForbidden, "checkout your information")
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logger.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return nil, e.Error(http.StatusInternalServerError, "get bucket failed")
	}
	path = filepath.Clean(path)

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			logger.FromContext(ctx).Errorf("not find bucket %s failed.Error:%v", bucket.Bucket, err)
			return nil, e.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
		}
		logger.FromContext(ctx).Errorf("find path(%s) failed.Error:%v", path, err)
		return nil, e.Error(http.StatusInternalServerError, "get bucket failed")
	}
	var info = &Info{
		LastModified: fileInfo.ModTime(),
		Name:         fileInfo.Name(),
	}
	if info.Size, info.Count, err = util.DirSize(path); err != nil {
		logger.FromContext(ctx).Errorf("lookup path %s failed.Error:%v", path, err)
		return nil, e.Error(http.StatusInternalServerError, "find bucket's all file failed")
	}
	return info, nil
}

// DeleteBucket 删除桶
func DeleteBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()

	bucket := &db.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?", bucketName, token.Domain).Find(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return e.Error(http.StatusNotFound, fmt.Sprintf("not found bucket %s", bucket.Bucket))
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
		return e.Error(http.StatusInternalServerError, "get bucket failed")
	}
	path = filepath.Clean(path)
	if err = os.RemoveAll(path); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("remove %s failed.Error:%v", path, err)
		return e.Error(http.StatusInternalServerError, "delete bucket failed")
	}
	return nil
}
