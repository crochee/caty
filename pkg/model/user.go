// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	"obs/pkg/db"
	"obs/pkg/log"
)

type User struct {
	ID uint64 `json:"id" gorm:"primary_key:id"`

	User string `gorm:"column:user;type:varchar(50);not null"`

	Domain string `gorm:"column:domain;type:varchar(50);not null"`

	Nick     string `gorm:"type:varchar(50);not null;column:nick"`
	PassWord string `gorm:"type:varchar(20);not null;column:pass_word"`

	Permission string `gorm:"type:text;not null;column:permission"`

	Deleted db.Deleted
	db.Base
}

func DeleteUser() {
	u := new(User)
	if err := db.Client().Model(u).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(u).Error; err != nil {
		log.Warn(err.Error())
	}
}
