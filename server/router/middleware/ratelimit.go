package router

import (
	"strconv"

	"github.com/AirGo-Official/AirGo/constant"
	"github.com/AirGo-Official/AirGo/global"
	"github.com/AirGo-Official/AirGo/utils/response"
	"github.com/gin-gonic/gin"
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
