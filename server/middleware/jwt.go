package middleware

import (
	"AirGo/global"
	"AirGo/utils/jwt_plugin"
	"AirGo/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		var token string
		token = c.GetHeader("Authorization")
		//fmt.Println("token", token)
		//判断
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
		claims, err := jwt_plugin.ParseTokenHs256(token, global.Server.JWT.SigningKey)
		if err != nil { //token过期，或其他解析错误
			response.Result(response.TOKENERROR, err.Error(), nil, c)
			c.Abort()
			return
		}
		//log.Println("token解析后", claims)
		//设置user id
		c.Set("uID", claims.UserID)
		c.Set("uName", claims.UserName)
		c.Next()
	}

}
