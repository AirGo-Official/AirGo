package middleware

import (
	"AirGo/global"
	"AirGo/utils/format_plugin"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logrus.Warn("Recovery middleware panic:", format_plugin.ErrorToString(err))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}

}
