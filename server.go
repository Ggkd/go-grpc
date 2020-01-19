package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"github.com/Ggkd/v1/proto"
	"runtime/debug"
)

//type TestService struct {
//
//}
//
////实现接口
//func (ts *TestService) Test (ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error){
//	return &proto.TestResponse{
//		Response:             req.GetRequest() + " test server",
//	}, nil
//}


//日志截取器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

//异常保护截取器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}


func main() {
	//// 创建 gRPC Server 对象
	//grpcServer := grpc.NewServer()
	//// 将TestService 注册到注册中心
	//proto.RegisterTestServiceServer(grpcServer, &TestService{})
	//// 创建listen，监听端口
	//listen, err := net.Listen("tcp", "localhost:8886")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// grpcServer 开始listen
	//err = grpcServer.Serve(listen)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}




	//	 使用截取器

	// 创建 gRPC Server 对象 ， 并添加截取器
	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(LoggingInterceptor, RecoveryInterceptor))
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