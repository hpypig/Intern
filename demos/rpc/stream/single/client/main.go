package main

import (
    "client/pb"
    "context"
    "flag"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "io"
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
    runLotsOfReplies(c)
}
func runLotsOfReplies(c pb.GreeterClient) {
    // server端流式RPC
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
    if err != nil {
        log.Fatalf("c.LotsOfReplies failed, err: %v", err)
    }
    for {
        // 接收服务端返回的流式数据，当收到io.EOF或错误时退出
        res, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("c.LotsOfReplies failed, err: %v", err)
        }
        log.Printf("got reply: %q\n", res.GetReply())
    }
}
