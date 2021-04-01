// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package user

type Domain struct {
	Nick string `json:"nick" uri:"nick" form:"nick" binding:"required"`
	LoginInfo
}

type LoginInfo struct {
	Email    string `json:"email" uri:"email" form:"email" binding:"required"`
	PassWord string `json:"pass_word" uri:"pass_word" form:"pass_word" binding:"required"`
}

type ModifyInfo struct {
	Email       string `json:"email" uri:"email" form:"email" binding:"required"`
	Nick        string `json:"nick" uri:"nick" form:"nick"`
	NewPassWord string `json:"new_pass_word" uri:"new_pass_word" form:"new_pass_word"`
	OldPassWord string `json:"old_pass_word" uri:"old_pass_word" form:"old_pass_word"`
}
