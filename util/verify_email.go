// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package util

import "regexp"

var regexpVerifyEmail = regexp.MustCompile(`^[A-Za-z\d]+([-_.][A-Za-z\d]+)*@([A-Za-z\d]+[-.])+[A-Za-z\d]{2,4}$`)

// VerifyEmail verify email
func VerifyEmail(email string) bool {
	return regexpVerifyEmail.MatchString(email)
}
