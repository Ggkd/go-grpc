package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"github.com/Ggkd/v1/proto"
)

type AuthenticationClient struct {
	User string
	Password string
}

func (a *AuthenticationClient) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user":a.User, "password":a.Password}, nil
}

func (a *AuthenticationClient) RequireTransportSecurity() bool {
	return false
}


func main() {
	auth := AuthenticationClient{
		User:     "luffy",
		Password: "123123",
	}
	//建立链接
	conn, err := grpc.Dial("localhost:8886", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
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