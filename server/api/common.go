package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
)

// gin.Context中获取user id
func GetUserIDFromGinContext(ctx *gin.Context) (int64, bool) {
	userID, ok := ctx.Get(constant.CTX_SET_USERID)
	return userID.(int64), ok
}

// gin.Context中获取user name
func GetUserNameFromGinContext(ctx *gin.Context) (string, bool) {
	userName, ok := ctx.Get(constant.CTX_SET_USERNAME)
	return userName.(string), ok
}

func EtagHandler(data any, ctx *gin.Context) {
	var md5, str string
	b, err := json.Marshal(data)
	if err != nil {
		ctx.AbortWithStatus(404)
		return
	}
	str = string(b)
	md5 = encrypt_plugin.Md5Encode(str, false)
	if md5 == ctx.Request.Header.Get("If-None-Match") {
		ctx.JSON(304, nil)
		return
	}
	ctx.Writer.Header().Set("Etag", md5)
	ctx.String(200, str)

}
