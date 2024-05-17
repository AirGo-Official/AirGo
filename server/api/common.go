package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"github.com/ppoonk/AirGo/utils/websocket_plugin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// gin.Context中获取user id
func GetUserIDFromGinContext(ctx *gin.Context) (int64, bool) {
	userID, ok := ctx.Get(constant.CTX_SET_USERID)
	if ok {
		return userID.(int64), ok
	} else {
		return 0, ok
	}
}

// gin.Context中获取user name
func GetUserNameFromGinContext(ctx *gin.Context) (string, bool) {
	userName, ok := ctx.Get(constant.CTX_SET_USERNAME)
	if ok {
		return userName.(string), ok
	} else {
		return "", ok
	}
}

func EtagHandler(data any, ctx *gin.Context) {
	var md5, str string
	b, err := json.Marshal(data)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	str = string(b)
	md5 = encrypt_plugin.Md5Encode(str, false)
	if md5 == ctx.Request.Header.Get("If-None-Match") {
		ctx.JSON(304, nil)
		return
	}
	ctx.Writer.Header().Set("Etag", md5)
	ctx.String(200, str)
}

// 升级websocket
func UpdateWebsocket(ctx *gin.Context, handler websocket_plugin.WebsocketMessageHandler) {
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
		global.Logrus.Error("websocket upgrade error:", err)
		response.Fail("websocket err:"+err.Error(), nil, ctx)
		return
	}
	client := &websocket_plugin.Client{
		ID:            strconv.FormatInt(uIDInt, 10),
		WsSocket:      conn,
		ClientChannel: make(chan *websocket_plugin.WsMessage, 2),
		ExpireTime:    20 * time.Second, //过期时间
	}
	defer client.Close()
	//2个协程去处理读写
	w := sync.WaitGroup{}
	w.Add(2)
	_ = global.GoroutinePool.Submit(func() {
		client.Read(&w, handler)
	})
	_ = global.GoroutinePool.Submit(func() {
		client.Write(&w)
	})

	w.Wait()
}

// Server-sent events
func SSE(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
