package main

import (
    "context"
    "fmt"
    "github.com/Shopify/sarama"
    "log"
    "strconv"
    "sync"
    "time"
)

/*

异步生产
1）生成 config 即可，不用配置
2）连接
3）defer Close
4）开启两个 goroutine 异步接收 successes、errors           这俩玩意儿究竟是啥啊？？
5）构造message，把其指针输入 producer.Input() 得到的 chan 即可，不用直接发送             这里应该是最体现异步的？？？
   另外，还要对这个输入操作，用 select 和 context 做超时控制
6）关闭连接，并等待返回结束
  producer.AsyncClose()
  wg.wait()

*/



// AsyncProducer 异步生产者
func AsyncProducer(topic string, limit int) {
    config := sarama.NewConfig()
    // 写这个的意思，大概就是，我打算在发送消息以后等待接收 successes、errors？？？？？
    // 异步生产者不建议把 Errors 和 Successes 都开启，一般开启 Errors 就行
    // 同步生产者就必须都开启，因为会同步返回发送成功或者失败
    config.Producer.Return.Errors = true // 设定是否需要返回错误信息
    config.Producer.Return.Successes = true // 设定是否需要返回成功信息
    producer, err := sarama.NewAsyncProducer([]string{"120.79.29.70:9092"}, config)
    if err != nil {
        log.Fatal("NewSyncProducer err:", err)
    }
    var (
        wg sync.WaitGroup
        enqueued, timeout, successes, errors int
    )
    // [!important] 异步生产者发送后必须把返回值从 Errors 或者 Successes 中读出来 不然会阻塞 sarama 内部处理逻辑 导致只能发出去一条消息
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 之前的demo里有 range 出现过阻塞，什么情况下会阻塞，这里不会阻塞？？？？？？？？？？？？？？？？！！！！！！！！！ 这种原理不注意可能发生死锁（或内存泄露？）
        for pm := range producer.Successes() {  // <-chan *ProducerMessage
            successes++
            fmt.Printf("%+v\n", pm)
        }
    }()
    wg.Add(1)
    go func() {
        defer wg.Done()
        for e := range producer.Errors() {
            log.Printf("[Producer] Errors：err:%v msg:%+v \n", e.Msg, e.Err)
            errors++
        }
    }()
    for i:=0; i<limit; i++ {
        str := strconv.Itoa(int(time.Now().UnixNano()))
        msg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(str)}
        // 异步发送只是写入内存了就返回了，并没有真正发送出去
        // sarama 库中用的是一个 channel 来接收，后台 goroutine 异步从该 channel 中取出消息并真正发送
        // select + ctx 做超时控制,防止阻塞 producer.Input() <- msg 也可能会阻塞
        ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
        select {
        case producer.Input() <- msg:
            enqueued++
        case <-ctx.Done():
            timeout++
        }
        cancel()
        if i%10000 == 0 && i != 0 {
            log.Printf("已发送消息数:%d 超时数:%d\n", i, timeout)
        }
    }
    // We are done
    producer.AsyncClose()
    // 事后处理结果，而不是刚发了等结果才能继续，这儿也体现异步
    wg.Wait()
    log.Printf("发送完毕 总发送条数:%d enqueued:%d timeout:%d successes: %d errors: %d\n", limit, enqueued, timeout, successes, errors)
}

