package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"github.com/Ggkd/v1/proto"
	"time"
)

type TestService struct {

}

//实现接口
func (ts *TestService) Test (ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error){
	for i:=0; i< 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "TestService.Test canceled")
		}
		time.Sleep(time.Second)
	}
	return &proto.TestResponse{
		Response:             req.GetRequest() + " test server",
	}, nil
}


func main() {
	// 创建 gRPC Server 对象
	grpcServer := grpc.NewServer()
	// 将TestService 注册到注册中心
	proto.RegisterTestServiceServer(grpcServer, &TestService{})
	// 创建listen，监听端口
	listen, err := net.Listen("tcp", "localhost:8886")
	if err != nil {
		fmt.Println(err)
		return
	}
	// grpcServer 开始listen
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}