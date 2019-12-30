package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/test/grpc_demo/proto"
)

func main() {
	//建立链接
	conn, err := grpc.Dial("localhost:8886", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}
	//打印响应
	fmt.Println("response------>>",resp.GetResponse())
}