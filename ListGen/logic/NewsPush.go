package logic

import (
    "ListGen/entities"
    "ListGen/midware"
    "ListGen/pb"
    "fmt"
    "google.golang.org/grpc"
    "log"
    "net"
    "sync"
)

// 暂时先写在这个文件夹，之后可能转移到rpc目录（包）下

type myServer struct {
    //UpdatedDataResChan chan entities.UpdatedDataResponse
    //Flag *SubFlag
    //FlagChan chan struct{}
    pb.UnimplementedNewsPusherServer
}

//type PushFlag struct {
//    flag bool
//    rw sync.RWMutex
//}
//func (pf PushFlag)isPush() bool {
//    pf.rw.RLock()
//    defer pf.rw.RUnlock()
//    return pf.flag
//
//}


func (s *myServer)ImportantNewsPush(in *pb.ImportantNewsRequest, stream pb.NewsPusher_ImportantNewsPushServer) error {  // 实现 proto service 里定义的方法
    //close(s.FlagChan) // 通知数据解析服务，此处 rpc 开始订阅数据
    //s.Flag.SetFlag(true) // 订阅
    // send
    ch := make(chan entities.UpdatedDataResponse, 10) // chan 在哪里关闭比较合适？？？
    go midware.GetUpdatedData(ch)

    for {
        fmt.Println("ImportantNewsPush - 等待更新数据...")
        updatedDataResponse := <- ch
        log.Printf("ImportantNewsPush-updatedDataResponse %+v\n", updatedDataResponse)
        var pbUpdatedNewsResponse pb.UpdatedNewsResponse
        CopyUpdatedDataResponse(&pbUpdatedNewsResponse, updatedDataResponse)
        log.Printf("ImportantNewsPush-pbUpdatedNewsResponse %+v\n", pbUpdatedNewsResponse)
        err := stream.Send(&pbUpdatedNewsResponse)
        if err != nil {
            log.Println("ImportantNewsPush err: ", err)
            return err
        }
    }
    //s.Flag.SetFlag(false) // 结束 rpc 时，停止订阅
    return nil
}





//func InitRPC(ch chan entities.UpdatedDataResponse, flagCh chan struct{}) {

func InitRPC() {
    // 初始化
    server := grpc.NewServer()
    pb.RegisterNewsPusherServer(server, &myServer{})

    lis, err := net.Listen("tcp",":8972")
    if err != nil {
        log.Println("PushNews err: ", err)
        return
    }
    err = server.Serve(lis)
    if err != nil {
        log.Println("PushNews err2: ", err)
        return
    }
}




func CopyUpdatedDataResponse(pbUpdatedData *pb.UpdatedNewsResponse, updatedData entities.UpdatedDataResponse) {
    pbUpdatedData.Op = updatedData.Op
    pbUpdatedData.Id = updatedData.Id
    pbUpdatedData.Data = &pb.ImportantNews{}
    log.Printf("CopyUpdatedDataResponse-updatedData.Data: %+v",updatedData.Data)
    CopyPbContent(pbUpdatedData.Data, updatedData.Data) // entities 里对应 NewsContent
}

func CopyPbContent(pbContent *pb.ImportantNews, content entities.NewsContent) {

    pbContent.Source = content.Source
    pbContent.Id = content.Id
    pbContent.Title = content.Title
    pbContent.Subtitle = content.Subtitle
    pbContent.Media = content.Media
    pbContent.Content = content.Content
    pbContent.PublishTime = content.PublishTime
    pbContent.Columns = content.Columns
    stocks := make([]*pb.Stock,len(content.Stocks))
    for i, stock := range content.Stocks {
        stocks[i] = &pb.Stock{}
        CopyPbStock(stocks[i], stock)
    }
    pbContent.Stocks  = stocks
    content.TxtType = pbContent.TxtType

}
func CopyPbStock( pbStock *pb.Stock, stock entities.Stock) {
    pbStock.Market = stock.Market
    pbStock.Code = stock.Code
    pbStock.Name = stock.Name

    tag := stock.Tag

    pbStock.Tag  = &pb.StockTag{
        Weight: tag.Weight,
        Emotion: tag.Emotion,
        EmotionWeight: tag.EmotionWeight,
    }
}

type SubFlag struct {
    flag bool
    rw sync.RWMutex
}
func (s SubFlag)OnPush() bool {
    s.rw.RLock()
    defer s.rw.RUnlock()
    return s.flag
}
func (s *SubFlag)SetFlag(flag bool) {
    s.rw.Lock()
    defer s.rw.Unlock()
    s.flag = flag
}

