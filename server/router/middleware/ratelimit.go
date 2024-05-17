package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/response"
	"strconv"
)

func RateLimitIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP() // localhost == ::1 时报错
		if ip == "::1" {
			ctx.Next()
		}
		if ok := global.RateLimit.IPRole.AllowVisit(ip); !ok {
			global.Logrus.Error(ip+"访问量超出,其剩余访问次数情况如下:", global.RateLimit.IPRole.RemainingVisits(ip))
			response.Response(constant.LIMITERROR, constant.ERROR_IP_LIMIT, nil, ctx)
			return
		}
		ctx.Next()
	}
}
func RateLimitVisit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uID, _ := ctx.Get(constant.CTX_SET_USERID)
		uIDStr := strconv.FormatInt(uID.(int64), 10)
		if ok := global.RateLimit.VisitRole.AllowVisit(uIDStr); !ok {
			global.Logrus.Error(uIDStr+"访问量超出,其剩余访问次数情况如下:", global.RateLimit.IPRole.RemainingVisits(uIDStr))
			response.Response(constant.LIMITERROR, constant.ERROR_USER_LIMIT, nil, ctx)
			return
		}
		ctx.Next()
	}
}
