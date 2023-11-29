package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/format_plugin"
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
