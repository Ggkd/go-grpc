package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type TestServer struct {

}

func (t *TestServer) Hello(request string, response *string) error {
	*response = "hello " + request
	return nil
}

func main() {
	//注册服务
	err := rpc.RegisterName("TestServer", new(TestServer))
	if err != nil {
		fmt.Println(err)
		return
	}
	//绑定端口
	listen,err := net.Listen("tcp", "localhost:8989")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		//建立连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 使用服务
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}