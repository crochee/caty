package e

// Code 规定组成部分为http状态码+5位错误码
type Code int

func (c Code) Error() string {
	if c < 10000000 || c >= 60000000 {
		return InvalidError
	}
	message, ok := codeMessageBox[c]
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
	// 00~100为服务级别错误码

	ErrInvalidParam        Code = 40010000
	ErrNotFound            Code = 40410001
	ErrInternalServerError Code = 50010002
)

var codeMessageBox = map[Code]string{
	ErrInvalidParam:        "请求参数不正确",
	ErrNotFound:            "资源不存在",
	ErrInternalServerError: "服务器内部错误",
}
