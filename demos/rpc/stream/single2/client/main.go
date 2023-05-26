package main

import (
    "client/pb"
    "context"
    "flag"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "log"
    "time"
)

const (
    defaultName = "world"
)

var (
    addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
    name = flag.String("name", defaultName, "Name to greet")
)

func main() {
    flag.Parse()
    // 连接到server端，此处禁用安全传输
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)
    runLotsOfGreeting(c)
}
func runLotsOfGreeting(c pb.GreeterClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    // 客户端流式RPC
    stream, err := c.LotsOfGreetings(ctx)
    if err != nil {
        log.Fatalf("c.LotsOfGreetings failed, err: %v", err)
    }
    names := []string{"七米", "q1mi", "沙河娜扎"}
    for _, name := range names {
        // 发送流式数据
        err := stream.Send(&pb.HelloRequest{Name: name})
        if err != nil {
            log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
        }
    }
    res, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("c.LotsOfGreetings failed: %v", err)
    }
    log.Printf("got reply: %v", res.GetReply())
}
