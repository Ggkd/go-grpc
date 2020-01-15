package main

import (
	"fmt"
	"net/rpc"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:8989")
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		var keyChange string
		err := client.Call("WatchServer.Watch", 5, &keyChange)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("keyChange:---", keyChange)
	}()

	err = client.Call("WatchServer.Set", [2]string{"name","luffy"}, new(struct{}))
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(7*time.Second)
}