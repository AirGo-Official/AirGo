package rpc_node_logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/rpc/server/node/rpc_node"

	"google.golang.org/grpc/metadata"
)

// HelloServer1 得有一个结构体，需要实现这个服务的全部方法,叫什么名字不重要
type HelloServer struct {
}

func (HelloServer) SayHello(ctx context.Context, request *rpc_node.HelloReq) (pd *rpc_node.HelloResp, err error) {

	//获取元数据信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传输token")
	}
	fmt.Println("获取元数据信息：", md)
	fmt.Println("入参：", request.Name)
	pd = new(rpc_node.HelloResp)
	pd.Name = "你好" + request.Name
	pd.Age = "ok" + request.Age
	return
}
