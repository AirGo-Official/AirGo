package public_api

import (
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/service/user_logic"
)

var (
	nodeService           *admin_logic.Node
	admin_customerService *admin_logic.CustomerService
	customerService       *user_logic.CustomerService
	shopService           *admin_logic.Shop
	userService           *user_logic.User
	payService            *user_logic.Pay
)
