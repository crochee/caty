// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import "github.com/crochee/lib/db"

type Account struct {
	ID   uint64 `json:"id" gorm:"primary_key:id"`
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null;index:idx_name_deleted,unique;comment:用户名"`

	Deleted db.Deleted `json:"deleted" gorm:"not null;index:idx_name_deleted,unique;comment:软删除记录id"`
	db.Base
}

func (Account) TableName() string {
	return "account"
}
