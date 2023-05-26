package main

import (
    "client/pb"
    "context"
    "flag"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "io"
    "log"
)

const (
    defaultName = "world"
)

var (
    addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
    name = flag.String("name", defaultName, "Name to greet")
)

/*
第一步：定义 service、function、message；生成相应 go 文件
第二步：连接、获取client、调用函数获取stream、等待接收
    conn,_ := grpc.Dial("host:port", grpc.WithTransportCredentials(insecure.NewCredentials()))
    client := pb.NewGreeterClient(conn) // Greeter 是自己在 proto 里定义的服务名
    stream,_ := client.Function(ctx, &pb.MessageName{})
    for { stream.Recv() }
*/


func main() {
    flag.Parse()
    // 连接到server端，此处禁用安全传输
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)
    //runLotsOfReplies(c)
    runUpdatedNewsReplies(c)
}
func runLotsOfReplies(c pb.GreeterClient) {
    // server端流式RPC
    //ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    //defer cancel()
    //stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
    stream, err := c.LotsOfReplies(context.Background(), &pb.HelloRequest{Name: *name})
    if err != nil {
        log.Fatalf("c.LotsOfReplies failed, err: %v", err)
    }
    for {
        // 接收服务端返回的流式数据，当收到io.EOF或错误时退出
        res, err := stream.Recv()
        if err == io.EOF {
            fmt.Println("eof")
            break
        }
        if err != nil {
            log.Fatalf("c.LotsOfReplies failed, err: %v", err)
        }
        log.Printf("got reply: %q\n", res.GetReply())
    }
}

func runUpdatedNewsReplies(c pb.GreeterClient) {
    //ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    //defer cancel()
    //stream, err := c.UpdatedNewsReplies(ctx, &pb.UpdatedNewsRequest{Nothing:"no params"})
    stream, err := c.UpdatedNewsReplies(context.Background(), &pb.UpdatedNewsRequest{Nothing:"no params"})
    if err != nil {
        log.Fatalf("c.LotsOfReplies failed, err: %v", err)
    }
    for {
        // 接收服务端返回的流式数据，当收到io.EOF或错误时退出
        res, err := stream.Recv()
        if err == io.EOF { // 什么情况会收 io.EOF ??????????
            fmt.Println("get eof")
            break
        }
        if err != nil {
            log.Fatalf("c.LotsOfReplies failed, err: %v", err)
        }
        log.Printf("got res: %+v\n", res)
        log.Printf("got stocks: %+v\n", res.GetStocks())
    }
}

