package middleware

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	//"AirGo/utils/casbin_plugin"
)

// Casbin 拦截器
func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		uID, ok := c.Get("uID")
		if !ok {
			response.Fail("接收参数错误", nil, c)
			c.Abort()
			return
		}
		uIdNew, _ := uID.(int64)
		//请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimSpace(path)
		//请求方法
		act := c.Request.Method
		// 获取用户的角色组
		var roleIds []int64
		global.DB.Model(&model.UserAndRole{}).Select("role_id").Where("user_id = ?", uIdNew).Find(&roleIds)
		//fmt.Println("获取用户的角色组:", roleIds)
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
