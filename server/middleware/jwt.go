package middleware

import (
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
		//fmt.Println("token", token)
		if token == "" {
			if token = c.GetHeader("Sec-WebSocket-Protocol"); token == "" {
				response.Fail("未携带token", nil, c)
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
			global.LocalCache.Delete(claims.UserName + "token") //删除过期token
			response.Result(response.TOKENERROR, err.Error(), nil, c)
			c.Abort()
			return
		}
		//设置user id
		c.Set(global.CtxSetUserID, claims.UserID)
		c.Set(global.CtxSetUserName, claims.UserName)
		c.Next()
	}

}
