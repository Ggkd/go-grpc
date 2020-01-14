package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/Ggkd/v1/proto"
	"time"
)

func main() {
	//建立链接
	conn, err := grpc.Dial("localhost:8886", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// context.WithDeadline：会返回最终上下文截止时间。第一个形参为父上下文，第二个形参为调整的截止时间。
	// 若父级时间早于子级时间，则以父级时间为准，否则以子级时间为最终截止时间
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	defer cancel()

	//创建客户端
	client := proto.NewTestServiceClient(conn)
	//发送请求
	resp, err := client.Test(ctx, &proto.TestRequest{
		Request:              "client request",
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("client err , deadline")
				return
			}
		}
		fmt.Println(err)
		return
	}
	//打印响应
	fmt.Println("response------>>",resp.GetResponse())
}