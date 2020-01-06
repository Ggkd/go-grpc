package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"project/test/grpc_demo/proto"
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


func main() {
	//***********使用证书**************//


	//生成根证书
	//生成的文件放在config目录下

	//	$ openssl genrsa -out ca.key 2048

	//	$ openssl req -new -x509 -days 3650 \
	//	-subj "/C=GB/L=China/O=gobook/CN=github.com" \
	//	-key ca.key -out ca.crt


	//   生成服务器端证书

	//	$ openssl req -new \
	//	-subj "/C=GB/L=China/O=server/CN=server.io" \
	//	-key server.key \
	//	-out server.csr

	//	$ openssl x509 -req -sha256 \
	//	-CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
	//	-in server.csr \
	//	-out server.crt


	//生成证书
	cred, err := tls.LoadX509KeyPair("./config/server.crt", "./config/server.key")
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
	// 将TestService 注册到注册中心
	proto.RegisterTestServiceServer(grpcServer, &TestService{})
	// 创建listen，监听端口
	listen, err := net.Listen("tcp", "localhost:8886")
	if err != nil {
		fmt.Println(err)
		return
	}
	// grpcServer 开始listen
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}