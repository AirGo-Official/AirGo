package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"strings"
)

func DomainAndAPI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.String()
		host := ctx.Request.Host

		//判断订阅域名,如果当前host在设置的订阅前缀里，则说明该host只是用来更新订阅的，只接受更新订阅的的请求，拒绝其他的请求
		if global.Server.Subscribe.SubscribeDomainBindRequest {
			if strings.Index(global.Server.Subscribe.BackendUrl, host) != -1 && strings.Index(path, "/api/public/sub") == -1 {
				ctx.Abort()
				return
			}
		}

		//是否启用静态资源 api (以 /api 开头的为业务api，其他的是静态资源api)
		if strings.Index(path, "/api") == -1 && !global.Server.Website.EnableAssetsApi {
			ctx.Abort()
			return
		}

		//是否启用swagger api
		if strings.Index(path, "/api/swagger") != -1 && !global.Server.Website.EnableSwaggerApi {
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
