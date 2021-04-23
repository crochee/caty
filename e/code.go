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

func (unknown) Chinese() string {
	return "发生未知错误"
}

type getTokenFail struct{}

func (g getTokenFail) String() string {
	return "OBS.00004"
}

func (g getTokenFail) Status() int {
	return http.StatusInternalServerError
}

func (g getTokenFail) English() string {
	return "Failed to obtain authentication information"
}

func (g getTokenFail) Chinese() string {
	return "获取鉴权信息失败"
}

type absPath struct{}

func (a absPath) String() string {
	return "OBS.00005"
}

func (a absPath) Status() int {
	return http.StatusInternalServerError
}

func (a absPath) English() string {
	return "Failed to get absolute path"
}

func (a absPath) Chinese() string {
	return "获取绝对路径失败"
}

type mkPath struct{}

func (m mkPath) String() string {
	return "OBS.00006"
}

func (m mkPath) Status() int {
	return http.StatusInternalServerError
}

func (m mkPath) English() string {
	return "Failed to generate path"
}

func (m mkPath) Chinese() string {
	return "生成路径失败"
}

type notFound struct{}

func (n notFound) String() string {
	return "OBS.00007"
}

func (n notFound) Status() int {
	return http.StatusNotFound
}

func (n notFound) English() string {
	return "Unable to find resource"
}

func (n notFound) Chinese() string {
	return "无法找到资源"
}

type operateDb struct{}

func (o operateDb) String() string {
	return "OBS.00008"
}

func (o operateDb) Status() int {
	return http.StatusInternalServerError
}

func (o operateDb) English() string {
	return "Failed to operate database"
}

func (o operateDb) Chinese() string {
	return "操作数据库失败"
}

type statisticsFile struct{}

func (s statisticsFile) String() string {
	return "OBS.00009"
}

func (s statisticsFile) Status() int {
	return http.StatusInternalServerError
}

func (s statisticsFile) English() string {
	return "Statistics file failed"
}

func (s statisticsFile) Chinese() string {
	return "统计文件失败"
}

type parseUrl struct{}

func (p parseUrl) String() string {
	return "OBS.00010"
}

func (p parseUrl) Status() int {
	return http.StatusBadRequest
}

func (p parseUrl) English() string {
	return "Failed to parse url"
}

func (p parseUrl) Chinese() string {
	return "解析请求链接失败"
}

type deleteBucket struct{}

func (d deleteBucket) String() string {
	return "OBS.00011"
}

func (d deleteBucket) Status() int {
	return http.StatusInternalServerError
}

func (d deleteBucket) English() string {
	return "Failed to delete bucket"
}

func (d deleteBucket) Chinese() string {
	return "删除桶失败"
}

type openFile struct{}

func (o openFile) String() string {
	return "OBS.00012"
}

func (o openFile) Status() int {
	return http.StatusInternalServerError
}

func (o openFile) English() string {
	return "Failed to open file"
}

func (o openFile) Chinese() string {
	return "打开文件失败"
}

type deleteFile struct{}

func (d deleteFile) String() string {
	return "OBS.00013"
}

func (d deleteFile) Status() int {
	return http.StatusInternalServerError
}

func (d deleteFile) English() string {
	return "Failed to delete file"
}

func (d deleteFile) Chinese() string {
	return "删除文件失败"
}

type generateToken struct{}

func (c generateToken) String() string {
	return "OBS.00014"
}

func (c generateToken) Status() int {
	return http.StatusInternalServerError
}

func (c generateToken) English() string {
	return "Failed to generate authentication information"
}

func (c generateToken) Chinese() string {
	return "生成鉴权信息失败"
}

type generateSign struct{}

func (c generateSign) String() string {
	return "OBS.00015"
}

func (c generateSign) Status() int {
	return http.StatusInternalServerError
}

func (c generateSign) English() string {
	return "Failed to generate signature information"
}

func (c generateSign) Chinese() string {
	return "生成签名信息失败"
}

type recovery struct{}

func (r recovery) String() string {
	return "OBS.00016"
}

func (r recovery) Status() int {
	return http.StatusInternalServerError
}

func (r recovery) English() string {
	return "A fatal error has occurred"
}

func (r recovery) Chinese() string {
	return "发生了致命错误"
}

type invalidEmail struct{}

func (i invalidEmail) String() string {
	return "OBS.00017"
}

func (i invalidEmail) Status() int {
	return http.StatusBadRequest
}

func (i invalidEmail) English() string {
	return "Invalid email address"
}

func (i invalidEmail) Chinese() string {
	return "无效的邮件地址"
}

type marshal struct{}

func (m marshal) String() string {
	return "OBS.00018"
}

func (m marshal) Status() int {
	return http.StatusInternalServerError
}

func (m marshal) English() string {
	return "Serialization failed"
}

func (m marshal) Chinese() string {
	return "序列化失败"
}

type unmarshal struct{}

func (u unmarshal) String() string {
	return "OBS.00019"
}

func (u unmarshal) Status() int {
	return http.StatusInternalServerError
}

func (u unmarshal) English() string {
	return "Deserialization failed"
}

func (u unmarshal) Chinese() string {
	return "反序列化失败"
}

type notAllow struct{}

func (n notAllow) String() string {
	return "OBS.00020"
}

func (n notAllow) Status() int {
	return http.StatusMethodNotAllowed
}

func (n notAllow) English() string {
	return "Method Not Allowed"
}

func (n notAllow) Chinese() string {
	return "不允许该方法"
}

type success struct{}

func (s success) String() string {
	return "OBS.00021"
}

func (s success) Status() int {
	return http.StatusOK
}

func (s success) English() string {
	return "success"
}

func (s success) Chinese() string {
	return "成功"
}
