// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"time"

	"gorm.io/gorm"
)

type BucketFile struct {
	File string `gorm:"primary_key:file"`

	Bucket string `gorm:"column:bucket;type:varchar(50);not null;"`

	Size    int64     `gorm:"column:size"`
	ModTime time.Time `gorm:"column:mod_time"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BucketFile) Delete() {
	tx := NewDB().Begin()
	defer tx.Rollback()
	if err := tx.Unscoped().Where("`deleted_at` IS NOT NULL").Delete(b).Error; err != nil {
		return
	}
	tx.Commit()
}
