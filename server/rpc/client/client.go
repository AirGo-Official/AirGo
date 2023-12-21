package main

import (
	"context"
	"fmt"
	"github.com/ppoonk/AirGo/rpc/server/node/rpc_node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// 自定义token认证
type CustomerTokenAuth struct {
}

// 获取元数据
func (c CustomerTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"key": "key",
	}, nil
}

// 是否开启传输安全 TLS
func (c CustomerTokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {

	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), //代表开启安全的选项
		grpc.WithPerRPCCredentials(new(CustomerTokenAuth)),       //添加自定义token验证
	}
	conn, err := grpc.Dial(":9988", opts...)
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect 连接失败", err))
	}
	defer conn.Close()
	// 初始化客户端
	client := rpc_node.NewHelloServiceClient(conn)
	result, err := client.SayHello(context.Background(), &rpc_node.HelloReq{
		Name: "小敏",
		Age:  "ok",
	})
	fmt.Println(result, err)
}
