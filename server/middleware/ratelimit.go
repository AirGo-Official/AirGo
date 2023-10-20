package middleware

import (
	"AirGo/global"
	"github.com/gin-gonic/gin"
	"strconv"
)

func RateLimitIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP() // localhost == ::1 时报错
		//fmt.Println("ClientIP:", ip)
		if ip == "::1" {
			ctx.Next()
		}
		if ok := global.RateLimit.IPRole.AllowVisit(ip); !ok {
			global.Logrus.Error(ip+"访问量超出,其剩余访问次数情况如下:", global.RateLimit.IPRole.RemainingVisits(ip))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
func RateLimitVisit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uID, _ := ctx.Get("uID")
		uIDStr := strconv.FormatInt(uID.(int64), 10)
		if ok := global.RateLimit.VisitRole.AllowVisit(uIDStr); !ok {
			global.Logrus.Error("访问量超出,其剩余访问次数情况如下:", global.RateLimit.VisitRole.RemainingVisits(uIDStr))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
