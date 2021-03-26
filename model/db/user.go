// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package db

import "time"

type User struct {
	User string `gorm:"primary_key:user;type:varchar(50);not null"`

	Domain string `gorm:"column:domain;type:varchar(50);not null"`

	Nick     string `gorm:"type:varchar(50);not null;column:nick"`
	PassWord string `gorm:"type:varchar(20);not null;column:pass_word"`

	Permission string `gorm:"type:text;not null;column:permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
