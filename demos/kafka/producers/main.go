package main

/*
ref: https://www.lixueduan.com/posts/kafka/05-quick-start/
     https://github.com/lixd/kafka-go-example
*/

/*
问题：
怎么检查子 goroutine 是否退出了？
主退出，子并不会退出。main退出，所有都退出。这是什么原理？
*/

func main() {
    SyncProducer("test_topic",2)
    //AsyncProducer("test_topic",2)
}


// 李文周----
//func main() {
//    // 指定 ack 方式（一致性强弱）
//    // 设置分区函数？？？？
//    // 成功交付的消息将在 success channel 返回
//    config := sarama.NewConfig()
//    config.Producer.RequiredAcks = sarama.WaitForAll
//    config.Producer.Partitioner = sarama.NewRandomPartitioner
//    config.Producer.Return.Successes = true
//
//    // 构造消息
//    msg := &sarama.ProducerMessage{}
//    msg.Topic = "web_log"  // topic 可能相当于队列？？？
//    msg.Value = sarama.StringEncoder("this is a test log")
//
//    // 连接kafka
//    client, err := sarama.NewSyncProducer([]string{"120.79.29.70:9092"}, config)
//    if err != nil {
//        log.Println("producer closed, err:", err)
//        return
//    }
//    defer client.Close()
//    // 发送消息         返回的是什么？有什么用？
//    pid, offset, err := client.SendMessage(msg) //
//    if err != nil {
//        log.Println(err)
//        return
//    }
//    fmt.Printf("pid:%v offset:%v\n", pid, offset)
//}




















































