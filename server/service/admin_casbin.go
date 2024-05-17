package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"strconv"
	"strings"
)

func Show(data any) {
	b, _ := json.Marshal(data)
	fmt.Println(string(b))
}

type Casbin struct{}

var AdminCasbinSvc *Casbin

func (c *Casbin) NewSyncedCachedEnforcer() (*casbin.SyncedCachedEnforcer, error) {
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
		return nil, err
	}
	a, _ := gormadapter.NewAdapterByDB(global.DB)
	syncedCachedEnforcer, _ := casbin.NewSyncedCachedEnforcer(m, a)
	syncedCachedEnforcer.SetExpireTime(60 * 60)
	_ = syncedCachedEnforcer.LoadPolicy()
	return syncedCachedEnforcer, nil
}

func (c *Casbin) UpdateCasbinPolicy(casbinInfoParams *model.CasbinInfo) error {
	roleID := strconv.FormatInt(casbinInfoParams.RoleID, 10)
	var list []gormadapter.CasbinRule
	var data []string
	for k, _ := range casbinInfoParams.CasbinItems {
		data = append(data, casbinInfoParams.CasbinItems[k].Path)
	}
	err := global.DB.Model(&gormadapter.CasbinRule{}).Where("v0 = 1 and v1 in (?)", data).Find(&list).Error
	if err != nil {
		return err
	}
	c.ClearCasbin(0, roleID)
	rules := [][]string{}
	for _, v := range list {
		rules = append(rules, []string{roleID, v.V1, v.V2})
	}
	success, err := global.Casbin.AddPolicies(rules)
	if !success {
		return errors.New("casbin添加失败")
	}
	err = global.Casbin.InvalidateCache()
	if err != nil {
		return err
	}
	//重新加载全局casbin
	res, err := c.NewSyncedCachedEnforcer()
	if err != nil {
		return err
	}
	global.Casbin = res
	return nil
}
func (c *Casbin) GetPolicyByRoleID(casbinInfo *model.CasbinInfo) *model.CasbinInfo {
	roleID := strconv.FormatInt(casbinInfo.RoleID, 10)
	list := global.Casbin.GetFilteredPolicy(0, roleID)
	//fmt.Println("GetPolicyByRoleID list:", list)
	for _, v := range list {
		casbinInfo.CasbinItems = append(casbinInfo.CasbinItems, model.CasbinItem{
			Path:   v[1],
			Method: v[2],
		})
	}
	return casbinInfo
}
func (c *Casbin) GetAllPolicy() *model.CasbinInfo {
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
func (c *Casbin) ClearCasbin(v int, p ...string) bool {
	success, _ := global.Casbin.RemoveFilteredPolicy(v, p...)
	return success
}
func (c *Casbin) GetUserAllRoutesByUserID(uid int64) (string, error) {
	var roleIDs []int64
	err := global.DB.Model(&model.UserAndRole{}).Select("role_id").Where("user_id = ?", uid).Find(&roleIDs).Error
	if err != nil {
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
	byte, err := json.Marshal(routesMap)
	if err != nil {
		return "", err
	}
	return string(byte), nil
}
