package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"log"
	"project/test/grpc_demo/proto"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

func main() {
	ZIPKIN_HTTP_ENDPOINT := "http://192.168.153.6:9411/api/v2/spans"
	//set up a span reporter
	reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	//create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("TestService", "localhost:8886")
	if err != nil {
		fmt.Println("endpoint---", err)
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

	//建立链接
	conn, err := grpc.Dial("localhost:8886", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		fmt.Println("conn---",err)
		return
	}
	defer conn.Close()
	//创建客户端
	client := proto.NewTestServiceClient(conn)
	//发送请求
	resp, err := client.Test(context.Background(), &proto.TestRequest{
		Request:              "client request",
	})
	if err != nil {
		fmt.Println("resp---", err)
		return
	}
	//打印响应
	fmt.Println("response------>>",resp.GetResponse())
}