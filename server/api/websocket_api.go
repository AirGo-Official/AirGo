package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
	"github.com/ppoonk/AirGo/utils/websocket_plugin"
	"net/http"
	"strconv"
	"time"
)

// websocket
func WebSocketMsg(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//后端带token响应，否则前端接收不到数据
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	uIDInt, ok := GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header"
		//nginx:
		//proxy_set_header Upgrade $http_upgrade;
		//proxy_set_header Connection upgrade;
		//proxy_set_header X-Real-IP $remote_addr;
		global.Logrus.Error("websocket upgrade error:", err)
		response.Fail("websocket err:"+err.Error(), nil, ctx)
		return
	}
	//defer conn.Close()
	client := &websocket_plugin.Client{
		ID: strconv.FormatInt(uIDInt, 10),
		//ID:            ctx.ClientIP(),
		WsSocket:      conn,
		ClientChannel: make(chan []byte),
		ExpireTime:    30 * time.Second, //过期时间
		QuitChanel:    make(chan bool),
	}
	global.WsManager.OnlineChannel <- client
	global.GoroutinePool.Submit(func() {
		client.Read(global.WsManager, service.GetNodesStatus)
	})
	global.GoroutinePool.Submit(func() {
		client.Write(global.WsManager)
	})
}
