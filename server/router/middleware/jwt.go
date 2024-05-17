package router

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/jwt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		var token string
		token = c.GetHeader("Authorization")
		if token == "" {
			if token = c.GetHeader("Sec-WebSocket-Protocol"); token == "" {
				response.Response(constant.TOKENERROR, constant.ERROR_NO_TOKEN_IN_THE_REQUEST, nil, c)
				c.Abort()
				return
			}
		}
		if strings.HasPrefix(token, "Bearer ") {
			//去掉bearer
			token = token[7:]
		}
		claims, err := jwt_plugin.ParseTokenHs256(token, global.Server.Security.JWT.SigningKey)
		if err != nil { //token过期，或其他解析错误
			response.Response(constant.TOKENERROR, err.Error(), nil, c)
			c.Abort()
			return
		}
		//设置user id
		c.Set(constant.CTX_SET_USERID, claims.UserID)
		c.Set(constant.CTX_SET_USERNAME, claims.UserName)
		c.Next()
	}

}
