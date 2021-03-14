// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

type DomainInfo struct {
	Email    string `gorm:"type:varchar(50);not null;unique_index;column:email"`
	Nick     string `gorm:"type:varchar(50);not null;column:nick"`
	PassWord string `gorm:"type:varchar(20);not null;column:pass_word"`
}
