// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	"github.com/crochee/lib/log"

	"cca/pkg/db"
)

type Account struct {
	ID uint64 `json:"id" gorm:"primary_key:id"`

	db.Base
}

func DeleteUser() {
	u := new(Account)
	if err := db.New().Model(u).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(u).Error; err != nil {
		log.Warn(err.Error())
	}
}
