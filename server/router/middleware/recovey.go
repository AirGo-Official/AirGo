package router

import (
	"runtime/debug"

	"github.com/AirGo-Official/AirGo/global"
	"github.com/AirGo-Official/AirGo/utils/format_plugin"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logrus.Warn("Recovery middleware panic:", format_plugin.ErrorToString(err), string(debug.Stack()))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
