package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	pb "github.com/leisurelicht/grpc-demo/protobuf"
	grpc "google.golang.org/grpc"
	"io"
	"net"
)

type authServer struct{}

func (*authServer) AuthLogin(stream pb.AUTH_AuthLoginServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Printf("收到的用户名：%s, 密码：%s.\n", req.Username, req.Password)
		c := make(chan string)
		go str2base64(c, req.Username, req.Password)
		for n := range c {
			resp := &pb.Response{
				Result: string(n),
			}
			stream.Send(resp)
		}
	}
	return nil

}

func newAuthServer() pb.AUTHServer {
	return &authServer{}
}

func str2base64(c chan string, username string, password string) {
	result := base64.StdEncoding.EncodeToString([]byte(username + password))
	fmt.Printf("Base64编码后的结果为: %s.\n", result)
	fmt.Printf("--------------------")
	c <- result
	close(c)
}

func main() {
	port := flag.Int("p", 12345, "服务运行端口")
	flag.Parse()

	fmt.Printf("认证服务启动, 运行端口为: %d", *port)
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAUTHServer(grpcServer, newAuthServer())
	grpcServer.Serve(conn)
}
