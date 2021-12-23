package resp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/afero/mem"
	"net/http"
	"strings"

	"github.com/crochee/lirity/e"
	"github.com/crochee/lirity/log"
	"github.com/gin-gonic/gin"
)

type ResponseCode struct {
	// 返回码
	// Required: true
	// Example: 20000000
	Code int `json:"code"`
	// 返回信息描述
	// Required: true
	// Example: success
	Msg string `json:"message"`
	// 具体描述信息
	Result string `json:"result"`
}

func (r *ResponseCode) Error() string {
	return fmt.Sprintf("{code:%d,message:%s,result:%s}", r.Code, r.Msg, r.Result)
}

// Error gin response with Code
func Error(ctx *gin.Context, code e.ErrorCode) {
	ctx.JSON(code.StatusCode(), ResponseCode{
		Code: code.Code(),
		Msg:  code.Message(),
	})
}

// A Wrapper provides context around another error.
type Wrapper interface {
	// Unwrap returns the next error in the error chain.
	// If there is no next error, Unwrap returns nil.
	Unwrap() error
}

// Errors gin Response with error
func Errors(ctx *gin.Context, err error) {
	log.FromContext(ctx.Request.Context()).Errorf("%+v", err)
	for err != nil {
		wrapper, ok := err.(Wrapper)
		if !ok {
			break
		}
		err = wrapper.Unwrap()
	}
	if err == nil {
		Error(ctx, e.ErrInternalServerError)
		return
	}
	var errorCode *e.ErrCode
	if errors.As(err, &errorCode) {
		Error(ctx, errorCode)
		return
	}
	Error(ctx, e.ErrInternalServerError)
}

// ErrorParam return error  with 404
func ErrorParam(ctx *gin.Context, err error) {
	log.FromContext(ctx.Request.Context()).Errorf("parse param failed.%+v", err)
	ctx.AbortWithStatusJSON(e.ErrInvalidParam.StatusCode(), ResponseCode{
		Code:   e.ErrInvalidParam.Code(),
		Msg:    e.ErrInvalidParam.Message(),
		Result: err.Error(),
	})
}

// Abort advance return error
func Abort(ctx *gin.Context, code e.ErrorCode, msg string) {
	ctx.AbortWithStatusJSON(code.StatusCode(), ResponseCode{
		Code:   code.Code(),
		Msg:    code.Message(),
		Result: msg,
	})
}

func Success(ctx *gin.Context, data ...interface{}) {
	if len(data) == 0 || data[0] == nil {
		ctx.Status(http.StatusNoContent)
		return
	}
	accept := ctx.Request.Header.Get("Accept")
	if strings.Contains(accept, "text/csv") && len(data) >= 2 {
		columnName, ok := data[1].(string)
		if !ok {
			columnName = "分布式服务"
		}
		file := mem.NewFileHandle(mem.CreateFile(columnName))
		err := json.NewEncoder(file).Encode(data[0])
		if err != nil {
			Errors(ctx, err)
			return
		}
		ctx.Render(http.StatusOK, NewCsvRender(file, ctx.Request))
		return
	}
	ctx.JSON(http.StatusOK, data[0])
}
