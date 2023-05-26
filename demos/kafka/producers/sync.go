package main

import (
    "github.com/Shopify/sarama"
    "log"
    "strconv"
    "time"
)

/*
同步生产
1）设置配置项
    同步配置：Success、Errors    不设置会怎样？
2）连接，得到 producer
3）defer Close
4）构造 message，使用 topic、key、value；value 需要编码？不编码会怎样？还有什么编码类型？
5）发送消息，返回 partition offset   貌似不能指定分区，分区是随机安排的？
*/

func SyncProducer(topic string, limit int) {
   config := sarama.NewConfig()

   //config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认

   // 轮循选择分区；也可以自己实现 Partitioner 接口来实现分区选择策略
   // API 提供 NewRandomPartitioner, NewHashPartitioner ...
   config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

   // 同步生产者必须同时开启 Return.Successes 和 Return.Errors
   // 因为同步生产者在发送之后就必须返回状态，所以需要两个都返回
   config.Producer.Return.Successes = true
   config.Producer.Return.Errors = true // 这个默认值就是 true 可以不用手动赋值
   // 同步生产者和异步生产者逻辑是一致的，Success或者Errors都是通过channel返回的，
   // 只是同步生产者封装了一层，等channel返回之后才返回给调用者
   // 具体见 sync_producer.go 文件72行 newSyncProducerFromAsyncProducer 方法
   // 内部启动了两个 goroutine 分别处理Success Channel 和 Errors Channel
   // 同步生产者内部就是封装的异步生产者
   // type syncProducer struct {
   // 	producer *asyncProducer
   // 	wg       sync.WaitGroup
   // }
   producer, err := sarama.NewSyncProducer([]string{"120.79.29.70:9092"}, config)
   if err != nil {
       log.Fatal("NewSyncProducer err:", err)
   }
   defer producer.Close()
   var successes, errors int
   for i := 0; i < limit; i++ {
       str := strconv.Itoa(int(time.Now().UnixNano()))
       msg := &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(str)}
       partition, offset, err := producer.SendMessage(msg) // 发送逻辑也是封装的异步发送逻辑，可以理解为将异步封装成了同步
       if err != nil {
           log.Printf("SendMessage:%d err:%v\n ", i, err)
           errors++
           continue
       }
       successes++
       // 第7条消息例子：2022/12/04 10:34:15 [Producer] partitionid: 0; offset:6, value: 1670121255401743700
       log.Printf("[Producer] partitionid: %d; offset:%d, value: %s\n", partition, offset, str)

   }
   log.Printf("发送完毕 总发送条数:%d successes: %d errors: %d\n", limit, successes, errors)
}
