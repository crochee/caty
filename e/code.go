// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/26

package e

import "net/http"

type Code interface {
	String() string
	Status() int
	English() string
	Chinese() string
}

type payload struct{}

func (payload) String() string {
	return "OBS.00001"
}

func (payload) Status() int {
	return http.StatusBadRequest
}

func (payload) English() string {
	return "Failed to parse request body"
}

func (payload) Chinese() string {
	return "解析请求体失败"
}

type forbidden struct{}

func (forbidden) String() string {
	return "OBS.00002"
}

func (forbidden) Status() int {
	return http.StatusForbidden
}

func (forbidden) English() string {
	return "Insufficient permission, operation forbidden"
}

func (forbidden) Chinese() string {
	return "权限不足，操作被禁止"
}

type unknown struct{}

func (unknown) String() string {
	return "OBS.00003"
}

func (unknown) Status() int {
	return http.StatusInternalServerError
}

func (unknown) English() string {
	return "An unknown error occurred"
}

func (u unknown) Chinese() string {
	return "发生未知错误"
}
