package main

import (
    "fmt"
    "google.golang.org/grpc"
    "net"
    "server/pb"
)

// 服务端流式rpc

// LotsOfReplies 返回使用多种语言打招呼
func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
    words := []string{
        "你好",
        "hello",
        "こんにちは",
        "안녕하세요",
    }

    for _, word := range words {
        data := &pb.HelloResponse{
            Reply: word + in.GetName(),
        }
        // 使用Send方法返回多个数据
        if err := stream.Send(data); err != nil {
            return err
        }
    }
    return nil
}

type server struct {
    pb.UnimplementedGreeterServer
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
