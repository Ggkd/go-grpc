package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"net"
	"github.com/Ggkd/v1/proto"
)

type AuthenticationServer struct {
	User string
	Password string
}

func (a *AuthenticationServer) GetAppKey() string {
	return "luffy"
}

func (a *AuthenticationServer) GetAppSecret() string {
	return "123123"
}

type TestService struct {
	auth *AuthenticationServer
}

//实现接口
func (ts *TestService) Test (ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error){
	if err := ts.auth.Auth(ctx); err != nil {
		return nil, err
	}
	return &proto.TestResponse{Response:"login success"}, nil
}

func (a *AuthenticationServer) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}

	var appid string
	var appkey string

	if val, ok := md["user"]; ok { appid = val[0] }
	if val, ok := md["password"]; ok { appkey = val[0] }

	if appid != a.GetAppKey() || appkey != a.GetAppSecret() {
		return grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
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