// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/25

package model

type AkSk struct {
	Ak string `json:"ak"  header:"ak" form:"ak" binding:"required"`
	Sk string `json:"sk" header:"sk" form:"sk" binding:"required"`
}
