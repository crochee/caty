// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"obs/config"
	"obs/e"
	"obs/internal"
	"obs/logger"
	"obs/model/db"
	"obs/service/business/tokenx"
)

// CreateBucket 创建桶
func CreateBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db.NewDB().Begin()
	defer tx.Rollback()
	bucket := &db.Bucket{
		Bucket: bucketName,
		Domain: token.Domain,
	}
	if err := tx.Create(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		logger.WithContext(ctx).Errorf("insert db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	path, err := filepath.Abs(fmt.Sprintf("%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logger.WithContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return e.New(e.GetAbsPathFail, "clear bucket failed")
	}
	if err = os.MkdirAll(filepath.Clean(path), os.ModePerm); err != nil {
		logger.WithContext(ctx).Errorf("create bucket failed.Error:%v", err)
		return e.New(e.MkPathFail, "create bucket failed")
	}
	tx.Commit()
	return nil
}

// HeadBucket 查询桶信息 Info
func HeadBucket(ctx context.Context, token *tokenx.Token, bucketName string) (*Info, error) {
	conn := db.NewDB()
	bucket := &db.Bucket{Bucket: bucketName}
	if err := conn.First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.Errorf(e.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logger.WithContext(ctx).Errorf("query db failed.Error:%v", err)
		return nil, e.New(e.OperateDbFail, err.Error())
	}
	if bucket.Domain != token.Domain {
		return nil, e.New(e.Forbidden, "checkout your information")
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logger.WithContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return nil, e.New(e.GetAbsPathFail, "clear bucket failed")
	}

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			logger.WithContext(ctx).Errorf("not find bucket %s failed.Error:%v", bucket.Bucket, err)
			return nil, e.Errorf(e.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logger.WithContext(ctx).Errorf("find path(%s) failed.Error:%v", path, err)
		return nil, e.New(e.OperateDbFail, err.Error())
	}
	var info = &Info{
		LastModified: fileInfo.ModTime(),
		Name:         fileInfo.Name(),
	}
	if info.Size, info.Count, err = internal.DirSize(path); err != nil {
		logger.WithContext(ctx).Errorf("lookup path %s failed.Error:%v", path, err)
		return nil, e.New(e.StatisticsFileFail, err.Error())
	}
	return info, nil
}

// DeleteBucket 删除桶
func DeleteBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db.NewDB().Begin()
	defer tx.Rollback()

	bucket := &db.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?", bucketName, token.Domain).
		First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.Errorf(e.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logger.WithContext(ctx).Errorf("query db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	if err := tx.Delete(bucket).Error; err != nil {
		logger.WithContext(ctx).Errorf("delete bucket failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logger.WithContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return e.New(e.GetAbsPathFail, "clear bucket failed")
	}
	path = filepath.Clean(path)
	if err = os.RemoveAll(path); err != nil {
		logger.WithContext(ctx).Errorf("remove %s failed.Error:%v", path, err)
		return e.New(e.DeleteBucketFail, "delete bucket failed")
	}
	tx.Commit()
	return nil
}
