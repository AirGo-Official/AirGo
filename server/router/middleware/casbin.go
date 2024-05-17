package router

import (
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Casbin 拦截器
func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		uID, ok := c.Get(constant.CTX_SET_USERID)
		if !ok {
			response.Fail("接收参数错误", nil, c)
			c.Abort()
			return
		}
		uIdInt, _ := uID.(int64)
		//请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimSpace(path)
		//请求方法
		act := c.Request.Method
		// 获取用户的角色组
		var roleIds []int64
		var err error
		//先判断cache中有无
		roleIdsCache, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_USER_ROLEIDS_BY_USERID, uIdInt))
		if ok {
			roleIds = roleIdsCache.([]int64)
		} else {
			roleIds, err = service.AdminRoleSvc.FindRoleIdsByuId(uIdInt)
			if err != nil {
				global.Logrus.Error("Casbin error:", err)
				c.Abort()
				return
			}
			if len(roleIds) == 0 {
				global.Logrus.Error("Casbin error:", "roleIds = 0")
				c.Abort()
				return
			}

		}
		//log.Println("roleIds:", roleIds)
		//log.Println("casbin sub:", strconv.FormatInt(v, 10))
		//log.Println("casbin obj:", obj)
		//log.Println("casbin act:", act)
		status := false
		for _, v := range roleIds {
			success, err := global.Casbin.Enforce(strconv.FormatInt(v, 10), obj, act) // 判断策略中是否存在
			if err != nil {
				global.Logrus.Error("权限casbin error:", err)
				response.Fail("权限错误"+err.Error(), nil, c)
				c.Abort()
				return
			}
			if success {
				status = true
				break
			}
		}
		if !status {
			global.Logrus.Error("权限不足:", obj)
			response.Fail("权限不足"+obj, nil, c)
			c.Abort()
			return
		}
		global.LocalCache.Set(fmt.Sprintf("%s%d",
			constant.CACHE_USER_ROLEIDS_BY_USERID, uIdInt),
			roleIds,
			constant.CACHE_CASBIN_ROLEIDS_TIMEOUT*time.Minute)
		//将角色组写入 gin.Context
		c.Set("roleIds", roleIds)
		c.Next()

	}
}
