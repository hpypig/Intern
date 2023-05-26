package main

import (
    "fmt"
    "google.golang.org/grpc"
    "log"
    "net"
    "server/pb"
    "strconv"
    "time"
)

type ImportantNews struct {
    Title string
    Media string
    Stocks []Stock
}
type Stock struct {
    Market string
    Code string
}

// 服务端流式rpc

func (s *server) UpdatedNewsReplies(in *pb.UpdatedNewsRequest, stream pb.Greeter_UpdatedNewsRepliesServer) error {

    fmt.Println("调用方参数：",in.GetNothing())
    stocks := []Stock{
        Stock{"s1","code1"},
        Stock{"s2","code2"},
    }
    news := ImportantNews{
        Title: "test title",
        Media: "test media",
        Stocks: stocks,
    }
    responseStocks := make([]*pb.Stock, len(stocks))

    for i, stock := range stocks {
        responseStocks[i] = &pb.Stock {
            Market: stock.Market,
            Code: stock.Code,
        }
    }
    for j:=0; j<3;j++ {
        for i:=0; i<len(stocks); i++ {

            data := &pb.UpdatedNewsResponse{
                Title: news.Title + strconv.Itoa(i),
                Media: news.Media + strconv.Itoa(i),
                Stocks: responseStocks,
            }
            if err := stream.Send(data); err != nil {
                log.Println("UpdatedNewsReplies err: ", err)
                return err
            }
        }
        time.Sleep(5*time.Second)
    }

    return nil
}

type server struct {
    pb.UnimplementedGreeterServer
}

/*
1.定义服务器
    type server struct {
        pb.UnimplementedGreeterServer
    }

2.创建服务器、注册服务
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})

3.创建套接字、服务器接收套接字
    lis, err := net.Listen("tcp", ":8972")
    err = s.Serve(lis)


*/

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
