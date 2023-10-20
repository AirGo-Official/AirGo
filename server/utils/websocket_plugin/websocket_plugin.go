package websocket_plugin

import (
	"AirGo/model"
	"AirGo/utils/logrus_plugin"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 消息体的定义
type WsMessage struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

// 客户端结构体定义
type Client struct {
	ID            string
	IpAddress     string
	WsSocket      *websocket.Conn
	ClientChannel chan []byte
	ExpireTime    time.Duration
	QuitChanel    chan bool
}

// 客户端管理
type ClientManager struct {
	Clients        map[string]*Client //记录在线用户
	MapLock        sync.RWMutex       //map读写锁
	Broadcast      chan []byte        //触发消息广播
	OnlineChannel  chan *Client       //用户登陆
	OfflineChannel chan *Client       //用户退出
}

func NewManager() *ClientManager {
	return &ClientManager{
		Clients:        make(map[string]*Client),
		Broadcast:      make(chan []byte, 10),
		OnlineChannel:  make(chan *Client, 10),
		OfflineChannel: make(chan *Client, 10),
	}
}

// 启动Manager，管理所有client，并进行相应handler（广播，下线等）
func (m *ClientManager) NewClientManager() {
	//该goroutine负责处理用户上线
	go func() {
		for {
			select {
			case oneClient := <-m.OnlineChannel:
				m.MapLock.Lock()
				m.Clients[oneClient.ID] = oneClient
				m.MapLock.Unlock()
			}
		}
	}()
	//该goroutine负责处理用户下线
	go func() {
		for {
			select {
			case oneClient := <-m.OfflineChannel:
				err := oneClient.WsSocket.Close() //关闭该client连接
				if err != nil {
					logFile, err1 := logrus_plugin.SetOutputFile()
					if err1 == nil {
						mw := io.MultiWriter(os.Stdout, logFile)
						log.SetOutput(mw)
						log.Println("oneClient.WsSocket.Close() error:", err)
						logFile.Close()
					}
					continue
				}
				oneClient.QuitChanel <- true   //通知write 关闭
				time.Sleep(time.Second)        //延时1秒关闭chanel
				close(oneClient.ClientChannel) //关闭client chanel
				close(oneClient.QuitChanel)    //关闭client quit chanel
				m.MapLock.Lock()
				delete(m.Clients, oneClient.ID) //deleted client
				m.MapLock.Unlock()
			}
		}
	}()
	//该goroutine负责广播
	go func() {
		for {
			select {
			case msg := <-m.Broadcast:
				for _, oneClient := range m.Clients {
					oneClient.ClientChannel <- msg
				}
			}
		}
	}()
}

// 读取消息,根据不同的type发送到对应的channel
func (c *Client) Read(manager *ClientManager, f func() *[]model.NodeStatus) {
	// 把当前客户端注销
	defer func() {
		if err := recover(); err != nil {
			logFile, err1 := logrus_plugin.SetOutputFile()
			if err1 == nil {
				mw := io.MultiWriter(os.Stdout, logFile)
				log.SetOutput(mw)
				log.Println("websocket read error:", err)
				logFile.Close()
			}
		}
		//_ = c.WsSocket.Close()
		//fmt.Println("defer close ws read")
		manager.OfflineChannel <- c //通知manager 下线该client
	}()
	for {
		_, data, err := c.WsSocket.ReadMessage() //阻塞
		if err != nil {
			fmt.Println("closed network connection，close ws read chanel:", err.Error())
			return
		}
		var wsmsg WsMessage
		err = json.Unmarshal(data, &wsmsg)
		if err != nil {
			fmt.Println("json.Unmarshal error:", err)
			continue
		}
		switch wsmsg.Type {
		case 8:
			return
		case 9:
			// 利用心跳监测
			resp, _ := json.Marshal(&WsMessage{Type: 10, Data: "pong"})
			c.ClientChannel <- resp
		case 1:
			// 推送TextMessage
			data := f()
			resp, _ := json.Marshal(&WsMessage{Type: 1, Data: data})
			c.ClientChannel <- resp
		}
	}
}

// 把对应消息写回客户端
func (c *Client) Write(manager *ClientManager) {
	defer func() {
		if err := recover(); err != nil {
			logFile, err1 := logrus_plugin.SetOutputFile()
			if err1 == nil {
				mw := io.MultiWriter(os.Stdout, logFile)
				log.SetOutput(mw)
				log.Println("websocket write error:", err)
				logFile.Close()
			}
		}
	}()
	for {
		select {
		case msg, ok := <-c.ClientChannel:
			if !ok {
				continue
			}
			err := c.WsSocket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("c.WsSocket.WriteMessage error:", err)
				return
			}
		case <-c.QuitChanel:
			return
		case <-time.After(c.ExpireTime):
			manager.OfflineChannel <- c
			return
		}
	}
}
