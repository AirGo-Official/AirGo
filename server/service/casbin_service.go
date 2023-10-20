package service

import (
	"AirGo/global"
	"AirGo/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"strconv"
	"strings"
)

func Casbin() *casbin.CachedEnforcer {
	//sub角色id，obj请求的路径，act请求的方法
	text := `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act 
	`
	//keyMatch2	一个URL 路径，例如 /alice_data/resource1,它返回一个布尔值表示 url 是否匹配
	m, err := casbinModel.NewModelFromString(text)
	if err != nil {
		global.Logrus.Error("casbin加载模型失败!", err.Error())
		return nil
	}
	a, _ := gormadapter.NewAdapterByDB(global.DB)
	cachedEnforcer, _ := casbin.NewCachedEnforcer(m, a) //生成casbin实施实例
	cachedEnforcer.SetExpireTime(60 * 60)
	err = cachedEnforcer.LoadPolicy()
	if err != nil {
		global.Logrus.Error("cachedEnforcer.LoadPolicy error:", err)
	}
	return cachedEnforcer
}

// 更新casbin权限
func UpdateCasbinPolicyOld(casbinInfo *model.CasbinInfo) error {
	roleIDInt := strconv.FormatInt(casbinInfo.RoleID, 10)
	ClearCasbin(0, roleIDInt)
	rules := [][]string{}
	for _, v := range casbinInfo.CasbinItems {
		rules = append(rules, []string{roleIDInt, v.Path, v.Method})
	}
	//fmt.Println("需要更新的casbin rules:", rules) //id 不能重复，负责casbin出错！
	//[[2 /user/login POST] [2 /user/getSub GET] [2 /casbin/getPolicyByRoleIds GET] [2 /casbin/updateCasbinPolicy POST] ]
	success, _ := global.Casbin.AddPolicies(rules)
	if !success {
		return errors.New("casbin添加失败")
	}
	err := global.Casbin.InvalidateCache()
	if err != nil {
		return err
	}
	//重新加载全局casbin
	global.Casbin = Casbin()
	return nil
}

// 更新casbin权限
func UpdateCasbinPolicy(casbinData *model.ChangeRoleCasbinReq) error {
	roleID := strconv.FormatInt(casbinData.RoleID, 10)
	var list []gormadapter.CasbinRule
	global.DB.Model(&gormadapter.CasbinRule{}).Where("v0 = 1 and v1 in (?)", casbinData.CasbinItems).Find(&list)

	ClearCasbin(0, roleID)
	rules := [][]string{}
	for _, v := range list {
		rules = append(rules, []string{roleID, v.V1, v.V2})
	}
	success, err := global.Casbin.AddPolicies(rules)
	if !success {
		global.Logrus.Error("casbin添加失败 err", success, err)
		//empty slice found
		return errors.New("casbin添加失败")
	}
	err = global.Casbin.InvalidateCache()
	if err != nil {
		return err
	}
	//重新加载全局casbin
	global.Casbin = Casbin()
	return nil
}

// API更新
func UpdateCasbinApi(oldPath, oldMethod, newPath, newMethod string) error {
	err := global.DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	err = global.Casbin.InvalidateCache()
	if err != nil {
		return err
	}
	return err
}

// 获取权限列表
func GetPolicyByRoleID(casbinInfo *model.CasbinInfo) *model.CasbinInfo {
	roleID := strconv.FormatInt(casbinInfo.RoleID, 10)
	list := global.Casbin.GetFilteredPolicy(0, roleID)
	for _, v := range list {
		casbinInfo.CasbinItems = append(casbinInfo.CasbinItems, model.CasbinItem{
			Path:   v[1],
			Method: v[2],
		})
	}
	return casbinInfo
}

// 获取全部权限
func GetAllPolicy() *model.CasbinInfo {
	var casbinInfo model.CasbinInfo
	list := global.Casbin.GetFilteredPolicy(0, "1") //超级管理员为全部api
	for _, v := range list {
		casbinInfo.CasbinItems = append(casbinInfo.CasbinItems, model.CasbinItem{
			Path:   v[1],
			Method: v[2],
		})
	}
	return &casbinInfo

}

// ClearCasbin
func ClearCasbin(v int, p ...string) bool {
	success, _ := global.Casbin.RemoveFilteredPolicy(v, p...)
	return success
}

// 获取用户全部路由 by user id
func GetUserAllRoutesByUserID(uid int64) (string, error) {
	var roleIDs []int64
	err := global.DB.Model(&model.UserAndRole{}).Select("role_id").Where("user_id = ?", uid).Find(&roleIDs).Error
	if err != nil {
		global.Logrus.Error(err.Error())
		return "", err
	}
	// ORM 进行去重查询，使用 Select 方法，将 DISTINCT 关键字嵌入
	var routes []string
	global.DB.Model(&gormadapter.CasbinRule{}).Select("DISTINCT v1").Where("v0 in ?", roleIDs).Find(&routes)
	// 使用/后面的字符串作为key
	routesMap := make(map[string]string, 0)
	for _, v := range routes {
		vv := v[strings.LastIndex(v, "/")+1:]
		routesMap[vv] = v
	}
	fmt.Println(len(routesMap))
	byte, err := json.Marshal(routesMap)
	if err != nil {
		return "", err
	}
	return string(byte), nil
}
