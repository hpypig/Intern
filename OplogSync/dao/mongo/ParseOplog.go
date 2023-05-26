package mongo

import (
    "fmt"
    "gopkg.in/mgo.v2/bson"
    "log"
    "oplogsync/dao/redis"
    "oplogsync/entities"
    "time"
)

/*
从 mongodb oplog 提取索引相关数据，然后存到 redis
    更新逻辑：
    每 5 秒更新一次，每次记录最新时间戳，下次获取这个时间戳之后的更新；断电以后怎么办？
    我在 redis 记录这个时间戳？redis 开启持久化？

    1）启动时：
    初始化更新时间，把启动时间前24小时作为更新时间？
    从 redis 读，如果redis为空，则设置为 24 小时之前的。
*/

func GetOplogAfter(lastTime int64, oplogChan chan<- entities.Oplog) {
    sess := session.Copy()

    //t, _ := time.Parse("2006-01-02 15:04:05","2022-12-11 16:10:57")

    if lastTime == 0 {
        //lastTime = time.Now().Add(-time.Hour * 72).Unix() // 1670733911  1670833119
        //lastTime = t.Unix()<<32
        lastTime = 1670746737-40*24*3600 // 1670746710
        fmt.Println("parse newLastTime: ", lastTime)
    }
    //fmt.Println("lastTime：", lastTime)

    // 从 oplog 查询大于 lastTime，op=i 的   要闻、个股、栏目
    //   1. 从 NewsContent 查 xx类型xx市场xxx股票 发布时间 id

    //fmt.Println("lastTime in query: ", lastTime+1)
    // 把某个时间以后的 content 索引全部更新到 redis
    // 为啥不直接用 64 位的时间？好像是为了和上面保持一致？

    var count int

    query := bson.M{"ts":bson.M{"$gte":bson.MongoTimestamp((lastTime+1)<<32)},"ns":"gf.NewsContent"}
    iter := sess.DB("local").C("oplog.rs").Find(query).Sort("$natural").Tail(5*time.Second)
    var oplog entities.Oplog
    for {
        for iter.Next(&oplog) {
            fmt.Println("新一轮iter")
            count++
            //fmt.Printf("ParseOplogAfter 资讯总数: %d, newsId: %s\n",count, oplog.O.Id.Hex()) // 424
            //fmt.Printf("oplog.ts:%d; _id: %s\n", oplog.Ts>>32, oplog.O.Id.Hex())

            oplogChan <- oplog

            // 我想每发一个 oplog 记录一次lastTime，但假如oplog解析失败，中途退出，那 lastTime 就记录错误了，会漏掉资讯
            //fmt.Printf("GetOplogAfter-oplog.Ts: %d\n",oplog.Ts)
            lastTime = oplog.Ts>>32
            _, err := redis.SetLastTime(lastTime)  // res
            if err != nil {
                log.Println("SetLastTime err: ", err)
            }
            fmt.Printf("GetOplogAfter - set lastime: %d\n",lastTime)
            oplog = entities.Oplog{}
        }
        var err error
        if err = iter.Err(); err != nil {
            log.Println("GetOplogAfter err: ", err)
            return // 为什么会出这个错误
            // 每次都在获取最后一秒数据的时候失败，再次拉取又可以了
            // wsarecv: A connection attempt failed because the connected party did not properly respond after a
            // period of time, or established connection failed because connected host has failed to respond.
            // 这个问题没有了，又变成
            // read tcp 192.168.1.103:8943->120.79.29.70:30002: wsarecv: A connection attempt failed because the connected party
            // did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
        }
        if iter.Timeout() {
            continue
        }

        lastTime, err = redis.GetLastTime()
        if err != nil {
            log.Println("69 parse GetLastTime err: ", err)
        } else {
            fmt.Println("71 last time: ", lastTime)
        }

        query = bson.M{"ts":bson.M{"$gte":bson.MongoTimestamp((lastTime+1)<<32)},"ns":"gf.NewsContent"}
        iter = sess.DB("local").C("oplog.rs").Find(query).Sort("$natural").Tail(5*time.Second)
    }
}




//-----
// 查一个索引更新一个索引；全部查出来后自己去分索引（感觉前者简单一点）

func ParseOplog() { // NewsType int, market string, code string
   sess := session.Copy()
   db := sess.DB("local")
   oplog := db.C("oplog.rs")
   var data []entities.Oplog
   //ts := 1670746736-40*24*3600
   // ts := 1670746736 // 30 个
   ts := 1670746737 // 30 个
   oplog.Find(bson.M{"ns":"gf.NewsContent","ts":bson.M{"$gt":bson.MongoTimestamp(ts<<32)}}).All(&data) // 这个是完整的 11 个
   fmt.Println("len data:",len(data))
   fmt.Printf("%v\n", data[0].Ts>>32)
   fmt.Printf("%+v\n", data[0])
   sess.Close()
}

func OplogTail() {
    sess := session.Copy()
    var oplog entities.Oplog
    //ts := 1670746736-40*24*3600
    // ts := 1670746736 //
    ts := 1670746737 // 10 个，不知道为啥第一个没显示
    // tail 方式 只有 10 个数据，第一个(xxx,1)的没取到，不知道为什么
    iter := sess.DB("local").C("oplog.rs").Find(bson.M{"ns":"gf.NewsContent","ts":bson.M{"$gt":bson.MongoTimestamp(ts<<32)}}).Tail(3*time.Second)
    var count int
    for {
        for iter.Next(&oplog) {
            count++
            fmt.Printf("oplog.o.id: %s, oplog.ts: %v\n, count: %d", oplog.O.Id.Hex(), oplog.Ts>>32, count)
        }
        if err := iter.Err(); err!= nil {
            fmt.Println("err: ", err)
        }
        if iter.Timeout() {
            continue
        }
        fmt.Println("next epoch")
        iter = sess.DB("local").C("oplog.rs").Find(bson.M{"ns":"gf.NewsContent","ts":bson.M{"$gt":bson.MongoTimestamp(ts<<32)}}).Tail(3*time.Second)
    }
    sess.Close()
}
// 问题是：为什么课题里的时间过滤没起作用
