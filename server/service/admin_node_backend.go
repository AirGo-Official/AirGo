package service

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
)

type NodeBackendService struct {
}
type NodeBackendServiceMessage struct {
	Title string
	Data  any
}
type NodeStatusMessage struct {
	CustomerServerIDs []int64
	NodeTrafficLog    *model.NodeTrafficLog
}
type NodeTrafficMessage struct {
	NodeTrafficLog *model.NodeTrafficLog
	AGUserTraffic  *model.AGUserTraffic
}
type UpdateCustomerTrafficLogMessage struct {
	CustomerServerIDs []int64
	UserTrafficLogMap map[int64]model.UserTrafficLog
}
type UpdateCustomerTrafficUsedMessage struct {
	CustomerServerIDs   []int64
	CustomerServiceList *[]model.CustomerService
}

func NewNodeBackendService() *NodeBackendService {
	return &NodeBackendService{}
}

var NodeBackendSvc *NodeBackendService

func InitNodeBackendSvc() {
	NodeBackendSvc = NewNodeBackendService()
	NodeBackendSvc.StartTask()
}
func (n *NodeBackendService) StartTask() {
	ch, err := global.Queue.Subscribe(constant.NODE_BACKEND_TASK, 10)
	if err != nil {
		global.Logrus.Error(err)
		return
	}
	go func() {
		for v := range ch {
			msg := v.(*NodeBackendServiceMessage)
			switch msg.Title {
			case constant.NODE_BACKEND_TASK_TITLE_NODE_STATUS:
				data := msg.Data.(*NodeStatusMessage)
				AdminNodeSvc.UpdateNodeStatus(data.CustomerServerIDs, data.NodeTrafficLog)

			case constant.NODE_BACKEND_TASK_TITLE_NODE_TRAFFIC:
				data := msg.Data.(*NodeTrafficMessage)
				AdminNodeSvc.UpdateNodeTraffic(data.NodeTrafficLog, data.AGUserTraffic)

			case constant.NODE_BACKEND_TASK_TITLE_UPDATE_CUSTOMER_TRAFFICLOG:
				data := msg.Data.(*UpdateCustomerTrafficLogMessage)
				_ = AdminCustomerServiceSvc.UpdateCustomerServiceTrafficLog(data.UserTrafficLogMap, data.CustomerServerIDs)

			case constant.NODE_BACKEND_TASK_TITLE_UPDATE_CUSTOMER_TRAFFICUSED:
				data := msg.Data.(*UpdateCustomerTrafficUsedMessage)
				_ = AdminCustomerServiceSvc.UpdateCustomerServiceTrafficUsed(data.CustomerServiceList, data.CustomerServerIDs)
		
			default:

			}
		}
	}()
}
