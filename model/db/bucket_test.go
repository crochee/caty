// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"testing"

	"obs/config"
)

func TestBucket_TableName(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	Setup()
	test := &Bucket{
		Domain: "123",
		Bucket: "bucket12",
	}
	//var list []*Bucket
	tx := NewDB().Begin()
	defer tx.Commit()
	// 级联插入
	result := tx.Create(test)
	if err := result.Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	t.Log(test)
}

func TestQuery(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	Setup()
	bucket := &Bucket{Bucket: "bucket12"}
	tx := NewDB().Begin()
	defer tx.Commit()
	// 级联插入
	if err := tx.Find(bucket).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	t.Log(bucket)
}
