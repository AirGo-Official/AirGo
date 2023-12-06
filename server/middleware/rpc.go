package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/rpc/rpctest/rpc/hello"
	"github.com/ppoonk/AirGo/rpc/rpctest/service"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func Rpc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var grpcServer = grpc.NewServer()
		server := service.HelloServer{}
		// 将server结构体注册为gRPC服务。
		hello.RegisterHelloServiceServer(grpcServer, &server)
		// 判断是否是grpc
		if ctx.Request.ProtoMajor == 2 &&
			strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			// 按grpc方式来请求
			ctx.Status(http.StatusOK) // 我本地测试这里不加状态码，客户端GRPC收到的是404状态码
			fmt.Println("处理grpc")
			grpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			// 不要再往下请求了,防止继续链式调用拦截器
			ctx.Abort()
			return
		}
		// 当作普通api
		ctx.Next()
	}
}
