package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rpc_node_logic "github.com/ppoonk/AirGo/rpc/server/node/logic"
	"github.com/ppoonk/AirGo/rpc/server/node/rpc_node"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func Rpc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 判断是否是grpc
		if ctx.Request.ProtoMajor == 2 && strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			ctx.Status(http.StatusOK)
			fmt.Println("处理grpc")
			//注册服务
			var grpcServer = grpc.NewServer()
			rpc_node.RegisterHelloServiceServer(grpcServer, &rpc_node_logic.HelloServer{})
			grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}
		// 当作普通api
		ctx.Next()
	}
}
