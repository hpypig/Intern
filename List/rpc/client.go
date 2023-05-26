package rpc

import (
    "context"
    "encoding/json"
    "github.com/hpypig/Intern/entities"
    "github.com/hpypig/Intern/router"
    "github.com/hpypig/Intern/rpc/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "log"
    "time"
)

func Init() {
    // 调用 listGen 持续推送--------------------------------------------------
    //ch := make(chan entities.UpdatedDataResponse, 10) // 接收 listgen 推送的更新内容
    broadcast := router.Hub.GetBroadcastCh()
    go UpdatedNewsPushRPC(broadcast)
}

func connect() (stream pb.NewsPusher_ImportantNewsPushClient) {
    conn, err := grpc.Dial("localhost:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Println("Rpc Connect: ", err)
        return
    } else {
        log.Println("connect success") // 奇怪啊，服务端都没开，怎么没有err呢？而是在下面调用时才发err
    }
    client := pb.NewNewsPusherClient(conn)

    var errCnt int
    for stream == nil {
        stream, err = client.ImportantNewsPush(context.Background(), &pb.ImportantNewsRequest{Placeholder: "nothing"})
        if err != nil {
            log.Println("runImportantNewsPush err: ",err)
            //return
        }
        errCnt++
        if errCnt == 10 {
            log.Println("UpdatedNewsPushRPC：rpc 多次连接失败")
            return
        }
        time.Sleep(time.Second*2)
    }
    return stream
}

// UpdatedNewsPushRPC 接收更新内容，输出到 ch
func UpdatedNewsPushRPC(ch chan<- []byte) {
//    conn, err := grpc.Dial("localhost:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
//    if err != nil {
//        log.Println("Rpc Connect: ", err)
//        return
//    } else {
//        log.Println("connect success")
//    }
//    defer func() {
//        log.Println("UpdatedNewsPushRPC-exit clientRPC")
//        conn.Close()
//    }()
//    client := pb.NewNewsPusherClient(conn)
//    stream, err := client.ImportantNewsPush(context.Background(), &pb.ImportantNewsRequest{Placeholder: "nothing"})
//    if err != nil {
//        log.Println("runImportantNewsPush err: ",err)
//        return
//    }
    stream := connect()

    var errCnt int
    for {
        if stream == nil {
            return
        }
        //fmt.Println("UpdatedNewsPushRPC - 等待 rpc 数据...")
        pbUpdatedData, err := stream.Recv()
        if err != nil {
           log.Println("UpdatedNewsPushRPC - Recv err: ", err)
           stream = connect()
           if errCnt == 10 {
               log.Println("UpdatedNewsPushRPC：rpc 多次接收失败")
               return
           }
           log.Println("UpdatedNewsPushRPC - errCnt: ", errCnt)
           continue
        }
        //log.Printf("UpdatedNewsPushRPC-pbUpdatedData: %+v\n", pbUpdatedData)
        var updatedDataResponse entities.UpdatedDataResponse
        CopyPbUpdatedData(&updatedDataResponse, pbUpdatedData)
        //log.Printf("UpdatedNewsPushRPC-updatedDataResponse: %+v\n", updatedDataResponse)
        bytes, err := json.Marshal(updatedDataResponse)
        if err != nil {
            log.Println("UpdatedNewsPushRPC--marshal err: ", err)
            return
        }
        ch <- bytes
    }

}

func CopyPbUpdatedData(updatedDataResponse *entities.UpdatedDataResponse, pbUpdatedData *pb.UpdatedNewsResponse) {
    updatedDataResponse.Op = pbUpdatedData.Op
    updatedDataResponse.Id = pbUpdatedData.Id
    CopyPbContent(&updatedDataResponse.Data, pbUpdatedData.Data)
}

func CopyPbContent(content *entities.NewsContent, pbContent *pb.ImportantNews) {

    content.Source = pbContent.Source
    content.Id = pbContent.Id
    content.Title = pbContent.Title
    content.Subtitle = pbContent.Subtitle
    content.Media = pbContent.Media
    content.Content = pbContent.Content
    content.PublishTime = pbContent.PublishTime
    content.Columns = pbContent.Columns
    stocks := make([]entities.Stock,len(pbContent.Stocks))
    for i, pbStock := range pbContent.Stocks {
        CopyPbStock(&stocks[i], pbStock) // 这一顿嵌套也太难受了呀
    }
    content.Stocks  = stocks
    content.TxtType = pbContent.TxtType

}
func CopyPbStock(stock *entities.Stock, pbStock *pb.Stock) {
    stock.Market = pbStock.Market
    stock.Code = pbStock.Code
    stock.Name = pbStock.Name

    pbTag := pbStock.Tag
    stock.Tag  = entities.StockTag{
        Weight: pbTag.Weight,
        Emotion: pbTag.Emotion,
        EmotionWeight: pbTag.EmotionWeight,
    }
}




