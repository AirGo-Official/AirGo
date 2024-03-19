package websocket_plugin

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 消息体的定义
type WsMessage struct {
	Type int `json:"type"`
	Data any `json:"data"`
}

// 客户端结构体定义
type Client struct {
	ID            string
	UserID        int64
	WsSocket      *websocket.Conn
	ClientChannel chan *WsMessage
	ExpireTime    time.Duration
}

type WebsocketMessageHandler func(wsMessage *WsMessage, msgChannel chan<- *WsMessage)

// 读取消息
func (c *Client) Read(w *sync.WaitGroup, f1 WebsocketMessageHandler) {
	defer w.Done()

	for {
		_ = c.WsSocket.SetReadDeadline(time.Now().Add(c.ExpireTime)) //每次读之前重置超时时间
		msgType, data, err := c.WsSocket.ReadMessage()               //阻塞。客户端关闭时会抛出err，然后Read退出
		if err != nil {
			fmt.Println("Read error:", err.Error())
			return
		}
		fmt.Println("msgType:", msgType, "data:", string(data))
		switch msgType {
		case websocket.CloseMessage:
			return
		case websocket.PingMessage:
			c.ClientChannel <- &WsMessage{Type: websocket.PongMessage, Data: "PONG"}
		default:
			if f1 != nil {
				f1(&WsMessage{
					Type: msgType,
					Data: data,
				}, c.ClientChannel)
			}
		}
	}
}

// 把对应消息写回客户端
func (c *Client) Write(w *sync.WaitGroup) {
	defer w.Done()
	for {
		select {
		case msg := <-c.ClientChannel:
			b, err := json.Marshal(msg.Data)
			if err != nil {
				continue
			}
			//_ = c.WsSocket.SetWriteDeadline(time.Now().Add(c.ExpireTime))
			_ = c.WsSocket.WriteMessage(msg.Type, b)
		case <-time.After(c.ExpireTime):
			return
		}
	}
}
func (c *Client) Close() {
	_ = c.WsSocket.Close()
	close(c.ClientChannel)
}
