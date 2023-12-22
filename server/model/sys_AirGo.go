package model

import (
	uuid "github.com/satori/go.uuid"
)

type AGNodeStatus struct {
	ID int64 `json:"id"`
	AGNodeStatusItem
}

// NodeStatus Node status
type AGNodeStatusItem struct {
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime uint64
}

type AGAccess struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Route string `json:"route"`
}
type AGUserInfo struct {
	ID             int64     `json:"id"`
	UUID           uuid.UUID `json:"uuid"`
	Passwd         string    `json:"passwd"`
	UserName       string    `json:"user_name"`
	NodeConnector  int64     `json:"node_connector"` //连接客户端数
	NodeSpeedLimit int64     `json:"node_speedLimit"`
}

type AGUserTraffic struct {
	ID          int64               `json:"id"`
	UserTraffic []AGUserTrafficItem `json:"user_traffic"`
}

type AGUserTrafficItem struct {
	UID      int64
	Email    string
	Upload   int64 //Byte
	Download int64 //Byte
}

type AGREALITYConfig struct {
	Dest             string
	ProxyProtocolVer uint64
	ServerNames      []string
	PrivateKey       string
	MinClientVer     string
	MaxClientVer     string
	MaxTimeDiff      uint64
	ShortIds         []string
}

type AGREALITYx25519 struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type AGOnlineUser struct {
	NodeID      int64
	UserNodeMap map[int64][]string //key:uid value:node ip array
}
