package user_api

import (
	"github.com/ppoonk/AirGo/service/user_logic"
)

var (
	orderService    *user_logic.Order
	shopService     *user_logic.Shop
	couponService   *user_logic.Coupon
	articleService  *user_logic.Article
	ticketService   *user_logic.Ticket
	userService     *user_logic.User
	customerService *user_logic.CustomerService
	roleService     *user_logic.Role
	menuService     *user_logic.Menu
	trafficService  *user_logic.Traffic
)
