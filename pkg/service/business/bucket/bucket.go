// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"cca/pkg/code"
	"cca/pkg/db"
	"cca/pkg/logx"
	"cca/pkg/model"
	"cca/pkg/service/business/tokenx"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"cca/config"

	"cca/internal"
)

// CreateBucket 创建桶
func CreateBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db.With(ctx).Begin()
	defer tx.Rollback()
	bucket := &model.Bucket{
		Bucket: bucketName,
		Domain: token.Domain,
	}
	if err := tx.Create(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.Wrap(code.ErrNotFound, "not found record")
		}
		logx.FromContext(ctx).Errorf("insert db failed.Error:%v", err)
		return code.Wrap(code.OperateDbFail, err.Error())
	}
	path, err := filepath.Abs(fmt.Sprintf("%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logx.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return code.Wrap(code.GetAbsPathFail, "clear bucket failed")
	}
	if err = os.MkdirAll(filepath.Clean(path), os.ModePerm); err != nil {
		logx.FromContext(ctx).Errorf("create bucket failed.Error:%v", err)
		return code.New(code.MkPathFail, "create bucket failed")
	}
	tx.Commit()
	return nil
}

// HeadBucket 查询桶信息 Info
func HeadBucket(ctx context.Context, token *tokenx.Token, bucketName string) (*Info, error) {
	conn := db2.NewDBWithContext(ctx)
	bucket := &model.Bucket{Bucket: bucketName}
	if err := conn.First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.Errorf(code.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logx.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return nil, code.New(code.OperateDbFail, err.Error())
	}
	if bucket.Domain != token.Domain {
		return nil, code.New(code.Forbidden, "checkout your information")
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logx.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return nil, code.New(code.GetAbsPathFail, "clear bucket failed")
	}

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			logx.FromContext(ctx).Errorf("not find bucket %s failed.Error:%v", bucket.Bucket, err)
			return nil, code.Errorf(code.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logx.FromContext(ctx).Errorf("find path(%s) failed.Error:%v", path, err)
		return nil, code.New(code.OperateDbFail, err.Error())
	}
	var info = &Info{
		LastModified: fileInfo.ModTime(),
		Name:         fileInfo.Name(),
	}
	if info.Size, info.Count, err = internal.DirSize(path); err != nil {
		logx.FromContext(ctx).Errorf("lookup path %s failed.Error:%v", path, err)
		return nil, code.New(code.StatisticsFileFail, err.Error())
	}
	return info, nil
}

// DeleteBucket 删除桶
func DeleteBucket(ctx context.Context, token *tokenx.Token, bucketName string) error {
	tx := db2.NewDBWithContext(ctx).Begin()
	defer tx.Rollback()

	bucket := &model.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?", bucketName, token.Domain).
		First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.Errorf(code.NotFound, "not found bucket %s", bucket.Bucket)
		}
		logx.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return code.New(code.OperateDbFail, err.Error())
	}
	if err := tx.Delete(bucket).Error; err != nil {
		logx.FromContext(ctx).Errorf("delete bucket failed.Error:%v", err)
		return code.New(code.OperateDbFail, err.Error())
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s",
		config.Cfg.ServiceConfig.ServiceInfo.StoragePath, bucket.Bucket))
	if err != nil {
		logx.FromContext(ctx).Errorf("get path abs failed.Error:%v", err)
		return code.New(code.GetAbsPathFail, "clear bucket failed")
	}
	path = filepath.Clean(path)
	if err = os.RemoveAll(path); err != nil {
		logx.FromContext(ctx).Errorf("remove %s failed.Error:%v", path, err)
		return code.New(code.DeleteBucketFail, "delete bucket failed")
	}
	tx.Commit()
	return nil
}
