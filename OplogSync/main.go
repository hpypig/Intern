package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/Shopify/sarama"
    "log"
    "oplogsync/dao/mongo"
    "oplogsync/dao/redis"
    "oplogsync/entities"
    "oplogsync/logic"
    "sync"
    "time"
)
const (
    Topic = "UpdatedData"
)

func main() {
    redis.Init()
    defer redis.Close()
    mongo.Init()
    defer mongo.Close()

    oplogChan := make(chan entities.Oplog, 10)
    lastTime, err := redis.GetLastTime()
    if err != nil {
        log.Println("~GetLastTime err: ", err)
    }
    log.Println("ListGen-lastTime: ", lastTime)
    go mongo.GetOplogAfter(lastTime, oplogChan)


    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    config.Producer.MaxMessageBytes = 10485880 // 避免消息过大被拦截

    producer,err := sarama.NewAsyncProducer([]string{"120.79.29.70:9092"},config)
    if err != nil {
        log.Println("main-NewAsyncProducer err: ", err)
        return
    }
    var successes, errors, enqueued int
    var wg sync.WaitGroup
    // 发送之前，开启接收 success 和 errors 的 goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        for _ = range producer.Successes() { // 我不知道这两个循环什么时候结束
            successes++
            //log.Printf("Successes() Offset: %+v\n", v.Offset)
        }
    }()
    wg.Add(1)
    go func() {
        defer wg.Done()
        for v := range producer.Errors() {
            errors++
            var temp entities.UpdatedDataResponse
            be ,ok  := (v.Msg.Value).(sarama.ByteEncoder)
            if ok {
                err = json.Unmarshal(be,&temp)
                if err != nil {
                    log.Println("be Unmarshal err: ", err)
                } else {
                    log.Printf("large message id: %v\n", temp.Id )
                }
            }

            log.Printf("large message size: %v\n", v.Msg.Value.Length())
            log.Printf("Errors(): %v\n", v.Error())

        }
    }()

    t := time.NewTimer(10*time.Second)
    for {
        var oplog entities.Oplog
        t.Reset(10*time.Second)
        fmt.Println("main - get oplog")
        select {
        case oplog = <- oplogChan:
        case <-t.C:
            //log.Printf("enqueued:%d", enqueued)  // suc err 要加锁才能访问！！！！！！！！！有更好的方法吗
            log.Printf("oplog timeout-enqueued:%d successes:%d errors:%d\n", enqueued, successes, errors)
        }

        updatedData := logic.GetUpdatedData(oplog)
        if updatedData == nil {
            //log.Println("**********************************************")
            continue
        }
        if updatedData.Op == "" || updatedData.Id=="" {
            continue
            //log.Println("op,id: ",updatedData.Op, updatedData.Id)
            //log.Println("enqueued: ",enqueued,"---------------------------------------------------------------------")
        }
        bytes, err := json.Marshal(updatedData) // 序列化 &struct 和 struct 有区别吗
        if err != nil {
            log.Println("main-Marshal err: ", err)
            return
        }
        mess := &sarama.ProducerMessage{Topic:Topic, Key: nil, Value: sarama.ByteEncoder(bytes)}
        ctx,_ := context.WithTimeout(context.Background(), 10 * time.Second)
        fmt.Println("main - mess 输入 kafka 通道...")
        select {
        case producer.Input() <- mess:
            enqueued++
        case <-ctx.Done():
            log.Printf("timeout updatedData: %+v\n",  updatedData) // 在传输该 data 时超时
            redis.SetLastTime(oplog.Ts>>32)
            wg.Wait()
            log.Printf("enqueued:%d successes:%d errors:%d\n", enqueued, successes, errors)
            return
        }
    }
}
