package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR      = 1
	SUCCESS    = 0
	TOKENERROR = 401 //前端判断401错误，token过期
	LIMITERROR = 408 //前端判断408错误，限流
)

// 序列化器
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`

	//Error  string      `json:"error"`
}

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}
func OK(message string, data interface{}, c *gin.Context) {
	Result(SUCCESS, message, data, c)
}

func Fail(message string, data interface{}, c *gin.Context) {
	Result(ERROR, message, data, c)
}
