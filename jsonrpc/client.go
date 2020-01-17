package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		fmt.Println(err)
		return
	}
	var response string
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err = client.Call("TestServer.Hello", "luffy", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("resp:::", response)
}