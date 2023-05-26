package main

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "net"
    "server/pb"
)
/*

ref: https://www.liwenzhou.com/posts/Go/gRPC/
protocol buffer: https://developers.google.com/protocol-buffers/docs/gotutorial
                 快速入门：https://www.tizi365.com/archives/367.html
                 https://www.liwenzhou.com/posts/Go/Protobuf3-language-guide-zh/

protoc: https://github.com/protocolbuffers/protobuf/releases

*/

type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func main() {
    // 监听本地的8972端口
    lis, err := net.Listen("tcp", ":8972")
    if err != nil {
        fmt.Printf("failed to listen: %v", err)
        return
    }
    s := grpc.NewServer()                  // 创建gRPC服务器
    pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
    // 启动服务
    err = s.Serve(lis)
    if err != nil {
        fmt.Printf("failed to serve: %v", err)
        return
    }
}
