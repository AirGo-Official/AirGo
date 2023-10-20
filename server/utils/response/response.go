package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR      = 1
	SUCCESS    = 0
	TOKENERROR = 401
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

// sspanel  /mod_mu/users 响应
type SSReP struct {
	Ret  int         `json:"ret"`
	Data interface{} `json:"data"`
}

func SSUsersOK(data interface{}, c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"ret":  1,
	// 	"data": data,
	// })
	c.JSON(http.StatusOK, SSReP{
		Ret:  1,
		Data: data,
	})
}
func SSUsersFail(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"ret":  0,
	// 	"data": nil,
	// })
	c.JSON(http.StatusOK, SSReP{
		Ret:  0,
		Data: nil,
	})
}
