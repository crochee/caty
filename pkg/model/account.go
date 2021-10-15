// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import "cca/pkg/db"

type Account struct {
	ID uint64 `json:"id" gorm:"primary_key:id"`

	db.Base
}

func (Account) TableName() string {
	return "account"
}
