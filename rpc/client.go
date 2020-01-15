package main

import (
	"fmt"
	"net/rpc"
)

//func main() {
//	client, err := rpc.Dial("tcp", "localhost:8989")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	var response string
//	err = client.Call("TestServer.Hello", "luffy", &response)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("resp:::", response)
//}


type TestServerClient struct {
	*rpc.Client
}

func (t *TestServerClient) Hello(request string, response *string) error {
	return t.Client.Call("TestServer.Hello", request, &response)
}

func DialTestServer(network, address string) (*TestServerClient, error) {
	client, err := rpc.Dial(network, address)
	return &TestServerClient{client}, err
}

//使用接口
func main() {
	client, err := DialTestServer("tcp", "localhost:8989")
	if err != nil {
		fmt.Println(err)
		return
	}
	var response string
	err = client.Hello("luffy", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("resp:::", response)
}