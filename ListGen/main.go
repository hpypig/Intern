package main

import (
    "ListGen/dao/mongo"
    "ListGen/dao/redis"
    "ListGen/logic"
    "fmt"
)

/*
解析oplog，传递给获取数据的服务
获取更新数据，传递给推送服务和 redis 更新服务

我想控制一个任务分发数据给两个chan，并且想随时控制开关，怎么写呢_发布订阅？
*/

func main() {
    redis.Init()
    defer redis.Close()
    mongo.Init()
    defer mongo.Close()


    //flagCh := make(chan struct{}) // 通知推送服务来订阅 updatedData 了 // 这种方法不行，多次 rpc 会 多次 close

    go logic.InitRPC()
    fmt.Println("rpc init")
    logic.UpdateRedis() // 可开多个，并发执行



    //go logic.PushNews(workChans[1], flagCH) // oplog chan；  grpc 等待获取 news // 不需要单独写 PushNews,grpc 部分已经实现了

    //oplogChan := make(chan entities.Oplog, 10)
    //go logic.GetUpdatedData(oplogChan, updatedDataChans, &subFlag) // 后期存入 kafka 而不是输入 chan
    ////go mongo.GetUpdatedData(oplogChan, subscribers) // 后期存入 kafka 而不是输入 chan
    //
    //
    //lastTime, err := redis.GetLastTime()
    //if err != nil {
    //    log.Println("~GetLastTime err: ", err)
    //}
    //fmt.Println("ListGen-lastTime: ", lastTime)
    ////workChans := []chan entities.Oplog{make(chan entities.Oplog,10), make(chan entities.Oplog,10)}
    //mongo.GetOplogAfter(lastTime, oplogChan) // 后期由 kafka 部分实现
    //// 等待请求，推送 important news 更新

}
