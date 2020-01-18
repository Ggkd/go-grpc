package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"project/test/grpc_demo/proto"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

type TestService struct {

}

//实现接口
func (ts *TestService) Test (ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error){
	return &proto.TestResponse{
		Response:             req.GetRequest() + " test server",
	}, nil
}


func main() {
	ZIPKIN_HTTP_ENDPOINT := "http://192.168.153.6:9411/api/v2/spans"
	//set up a span reporter
	reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	//create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("TestService", "localhost:8886")
	if err != nil {
		fmt.Println(err)
		return
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	// 创建 gRPC Server 对象
	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())))
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