package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"project/test/grpc_demo/proto"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial("localhost:8886", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建客户端
	client := proto.NewStreamServiceClient(conn)

	//// 1.基于服务的流
	//stream, err := client.BaseServer(context.Background(), &proto.StreamRequest{
	//	Name:                 "base on server --- client name.",
	//	Age:                  007,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for {
	//	resp, err := stream.Recv()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println("Name::", resp.GetName(), " age::", resp.GetAge())
	//}



	//// 2. 基于客户端的流
	//stream, err := client.BaseClient(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//for i:=0;i<10;i++ {
	//	err = stream.Send(&proto.StreamRequest{
	//		Name: "base on client --- client name.",
	//		Age:  int64(i),
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	time.Sleep(time.Second)
	//}
	//resp, err := stream.CloseAndRecv()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("Name::", resp.GetName(), " age::", resp.GetAge())



	//3.双向流
	stream, err := client.BaseDouble(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:=100;;i-- {
		err = stream.Send(&proto.StreamRequest{
			Name: "double --- client name.",
			Age:  int64(i),
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Name::", resp.GetName(), " age::", resp.GetAge())
	}
	//stream.CloseSend()
}
