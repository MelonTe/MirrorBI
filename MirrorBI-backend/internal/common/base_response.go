package common

import (
	"github.com/gin-gonic/gin"
	"mrbi/internal/ecode"
	"net/http"
)

// 统一响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data" swaggertype:"object"`
}

func BaseResponse(c *gin.Context, data interface{}, msg string, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
func Success(c *gin.Context, data interface{}) {
	BaseResponse(c, data, "", 0)
}

// 失败响应
func Error(c *gin.Context, code int) {
	BaseResponse(c, nil, ecode.GetErrMsg(code), code)
}
