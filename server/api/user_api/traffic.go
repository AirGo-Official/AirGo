package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/utils/response"
)

func GetSubTrafficList(ctx *gin.Context) {
	list, err := trafficService.GetSubTrafficList()
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", list, ctx)
}
