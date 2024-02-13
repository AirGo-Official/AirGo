package admin_api

import (
	"github.com/ppoonk/AirGo/service/admin_logic"
)

var (
	articleService  *admin_logic.Article
	menuService     *admin_logic.Menu
	orderService    *admin_logic.Order
	roleService     *admin_logic.Role
	systemService   *admin_logic.System
	shopService     *admin_logic.Shop
	userService     *admin_logic.User
	customerService *admin_logic.CustomerService
	ticketService   *admin_logic.Ticket
	casbinService   *admin_logic.Casbin
)
