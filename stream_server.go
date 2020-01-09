package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"github.com/Ggkd/v1/proto"
	"time"
)

type StreamServer struct {

}

func (ss *StreamServer) BaseServer(req *proto.StreamRequest,stream proto.StreamService_BaseServerServer) error {
	//基于服务的流
	for {
		err := stream.Send(&proto.StreamResponse{
			Name:                 req.GetName(),
			Age:                  req.GetAge(),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second)
	}
}

func (ss *StreamServer) BaseClient(stream proto.StreamService_BaseClientServer) error {
	//基于客户端的流
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.StreamResponse{
				Name:                 "base on client----server name.",
				Age:                  007,
			})
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Name::", recv.GetName(), " age::", recv.GetAge())
	}
}

func (ss *StreamServer) BaseDouble(stream proto.StreamService_BaseDoubleServer) error {
	for i:=1; ; i++ {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println("Name::", recv.GetName(), " age::", recv.GetAge())
		err = stream.Send(&proto.StreamResponse{
			Name:                 "double----server name.",
			Age:                  int64(i),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	// 创建grpc服务
	grpcServer := grpc.NewServer()
	//注册服务
	proto.RegisterStreamServiceServer(grpcServer, &StreamServer{})
	//创建listen
	listen, err := net.Listen("tcp", "localhost:8886")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 服务开始listen
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}