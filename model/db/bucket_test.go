// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"obs/config"
)

func TestBucket_TableName(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	if err := Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
	test := &Bucket{
		//Domain: "123",
		Bucket: "bucket12",
	}
	//var list []*Bucket
	tx := NewDB().Begin()
	defer tx.Commit()
	// 级联插入
	if err := tx.Create(test).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Delete(test).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	test.Bucket = ""
	if err := tx.Unscoped().Where("`deleted_at` IS NOT NULL").Delete(test).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	t.Log(test)
}

func TestMockBucket(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectExec("INSERT INTO obs_bucket").
		WithArgs("bucket12", "123", AnyTime{}, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))
	if err = db.Exec("INSERT INTO obs_bucket(bucket,domain,created_at,updated_at) VALUES (?,?,?,?)",
		"bucket12", "123", time.Now(), time.Now()).Error; err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQuery(t *testing.T) {
	config.InitConfig("../../conf/config.yml")
	if err := Setup(context.Background()); err != nil {
		t.Fatal(err)
	}
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
