// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import "time"

type BucketFile struct {
	ID       uint `gorm:"primary_key"`
	BucketId uint `gorm:"column:bucket_id"`

	File string `gorm:"column:file"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
