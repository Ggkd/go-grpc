package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net/http"
	"github.com/Ggkd/v1/proto"
	"strings"
)

//type TestService struct {
//
//}
//
////实现接口
//func (ts *TestService) Test (ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error){
//	return &proto.TestResponse{
//		Response:             req.GetRequest() + " test server",
//	}, nil
//}


func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" http server"))
	})

	return mux
}

func main() {
	//生成证书
	certFile := "./config/server.crt"
	keyFile := "./config/server.key"
	cred, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./config/ca.crt")
	if err != nil {
		fmt.Println(err)
		return
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cred},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	// 创建 gRPC Server 对象,使用证书
	grpcServer := grpc.NewServer(grpc.Creds(c))
	//创建一个新的 ServeMux，ServeMux 本质上是一个路由表。它默认实现了 ServeHTTP
	mux := GetHTTPServeMux()
	// 将TestService 注册到注册中心
	proto.RegisterTestServiceServer(grpcServer, &TestService{})
	err = http.ListenAndServeTLS("localhost:8886",
		certFile,
		keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}

			return
		}),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
