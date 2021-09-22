package e

// Code 规定组成部分为http状态码+5位错误码
type Code int

func (c Code) Error() string {
	if c < 10000000 || c >= 60000000 {
		return InvalidError
	}
	message, ok := codeZhMessageBox[c]
	if ok {
		return message
	}
	return UnDefineError
}

func (c Code) Code() int {
	return int(c)
}

func (c Code) StatusCode() int {
	return int(c) / 100000
}

const (
	UnDefineError = "未定义错误码"
	InvalidError  = "无效错误码"
)

const (
	ErrSuccess Code = 20000000
	// 00~99为服务级别错误码

	ErrInvalidParam        Code = 40010000
	ErrNotFound            Code = 40410001
	ErrInternalServerError Code = 50010002
	ErrMethodNotAllow      Code = 40510003
	ErrOperateDB           Code = 50010004

	// 100~200为用户验证类

	ErrInvalidEmail Code = 40510100
)

var codeZhMessageBox = map[Code]string{
	ErrSuccess: "成功",

	ErrInvalidParam:        "请求参数不正确",
	ErrNotFound:            "资源不存在",
	ErrMethodNotAllow:      "方法不允许",
	ErrInternalServerError: "服务器内部错误",
	ErrOperateDB:           "操作数据库失败",

	ErrInvalidEmail: "非法邮箱格式",
}
