// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	"github.com/crochee/lib/db"
	"github.com/crochee/lib/log"
)

type Domain struct {
	ID         uint64 `json:"id" gorm:"primary_key:id"`
	AccountID  string `json:"account_id" gorm:"column:account_id;type:varchar(50);not null;comment:主账号ID"`
	UserID     string `json:"user_id" gorm:"column:user_id;type:varchar(50);not null"`
	Nick       string `json:"nick" gorm:"type:varchar(50);not null;column:nick"`
	PassWord   string `json:"pass_word" gorm:"type:varchar(20);not null;column:pass_word"`
	Email      string `json:"email" gorm:"type:varchar(50);not null;unique_index:email"`
	Permission string `json:"permission" gorm:"type:text;not null;column:permission"`
	Verify     bool   `json:"verify" gorm:"column:verify"`

	Deleted db.Deleted `json:"deleted" gorm:""`
	db.Base
}

func DeleteDomain() {
	d := new(Domain)
	if err := New().Model(d).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(d).Error; err != nil {
		log.Warn(err.Error())
	}
}
