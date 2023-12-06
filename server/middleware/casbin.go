package middleware

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/ppoonk/AirGo/utils/casbin_plugin"
)

// Casbin 拦截器
func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		uID, ok := c.Get(global.CtxSetUserID)
		if !ok {
			response.Fail("接收参数错误", nil, c)
			c.Abort()
			return
		}
		uIdNew, _ := uID.(int64)
		uIdStr := fmt.Sprintf("%d", uIdNew)
		//请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimSpace(path)
		//请求方法
		act := c.Request.Method
		// 获取用户的角色组
		var roleIds []int64
		//先判断cache中有无
		roleIdsCache, ok := global.LocalCache.Get(uIdStr + global.UserRoleIds)
		if ok {
			roleIds = roleIdsCache.([]int64)
		} else {
			err := global.DB.Model(&model.UserAndRole{}).Select("role_id").Where("user_id = ?", uIdNew).Find(&roleIds).Error
			if err != nil {
				global.Logrus.Error("Casbin error:", err)
				c.Abort()
				return
			}
		}
		if len(roleIds) == 0 {
			global.Logrus.Error("Casbin error:", "roleIds = 0")
			c.Abort()
			return
		}
		global.LocalCache.Set(uIdStr+global.UserRoleIds, roleIds, 5*time.Minute) //超时
		status := false
		for _, v := range roleIds {
			success, err := global.Casbin.Enforce(strconv.FormatInt(v, 10), obj, act) // 判断策略中是否存在
			//log.Println("casbin sub:", strconv.FormatInt(v, 10))
			//log.Println("casbin obj:", obj)
			//log.Println("casbin act:", act)
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
		//将角色组写入 gin.Context
		c.Set("roleIds", roleIds)
		c.Next()

	}
}
