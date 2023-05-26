package main

import (
    "github.com/gin-gonic/gin"
    "github.com/hpypig/Intern/dao/mongo"
    "github.com/hpypig/Intern/dao/redis"
    "github.com/hpypig/Intern/router"
    "github.com/hpypig/Intern/rpc"
)
func main() {

    redis.Init()
    defer redis.Close()
    mongo.Init()
    defer mongo.Close()

    //c := redis.Client{}
    //RedisDataPrepare()
    //FindLatestNewsContent(0,"","")
    //mongo.IterTest()

    //go logic.ListGen() // 已经完成自动更新redis部分，还差 grpc 主动推送（插入、更新资讯）给 list 服务的部分

    gin.SetMode(gin.ReleaseMode)
    engine := gin.Default()
    router.Init(engine)
    rpc.Init()  // 这里有点bug，必须要求服务端先打开，否则rpc失效了；且rpc中途断开，不能重连。

    engine.Run(":8080")



}

