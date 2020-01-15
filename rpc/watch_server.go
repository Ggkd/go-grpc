package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type WatchServer struct {
	m      map[string]string
	filter map[string]func(key string)
	mu     sync.Mutex
}

func NewWatchServer() *WatchServer {
	return &WatchServer{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

func (w *WatchServer) Get(key string, reply *string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if v, ok :=w.m[key]; ok {
		*reply = v
		return nil
	}
	return fmt.Errorf("not found")
}

func (w *WatchServer) Set(kv [2]string, reply *struct{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	key, value := kv[0], kv[1]
	oldValue := w.m[key]
	if oldValue != value {
		fmt.Println(w.filter)
		for _, fn := range w.filter {
			fn(key)
		}
	}
	w.m[key] = value
	return nil
}

func (w *WatchServer) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	ch := make(chan string, 10)
	w.mu.Lock()
	w.filter[id] = func(key string) {
		ch <- key
	}
	defer w.mu.Unlock()
	select {
	case <- time.After(time.Duration(timeoutSecond)*time.Second):
		return fmt.Errorf("timeout")
	case key := <- ch:
		*keyChanged = key
		return nil
	}
}

func main() {
	//注册服务
	watchServer := NewWatchServer()
	watchServer.m["name"] = "zorro"
	err := rpc.RegisterName("WatchServer", watchServer)
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
	//建立连接
	conn, err := listen.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 启动服务
	rpc.ServeConn(conn)
}