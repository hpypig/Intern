package main

// ref: https://github.com/lixd/kafka-go-example/blob/main/consumer/group/group_consumer.go
//      https://www.lixueduan.com/posts/kafka/05-quick-start/

func main() {
    //SinglePartition("zyxz")
    SinglePartition("test_topic")
    //Partitions("zyxz")
    //OffsetManager("zyxz") // 这个运行了没反应，应该是因为我没有group，不知道这是个什么概念

    //go ConsumerGroup("test_topic","test_group","MyConsumerGroupHandler_name")
    //ConsumerGroup("test_topic","test_group","MyConsumerGroupHandler_name2")

}


/*
获取topic分区索引
对每个索引设置一个consumer
异步运行每个consumer
*/

//func main() {
//    consumer, err := sarama.NewConsumer([]string{"120.79.29.70:9092"}, nil)
//    if err != nil {
//        fmt.Printf("fail to start consumer, err:%v\n", err)
//        return
//    }
//    // 根据 topic 取到所有的分区
//    partitionList, err := consumer.Partitions("web_log")
//    if err != nil {
//        log.Println(err)
//        return
//    }
//    fmt.Println(partitionList)
//    for partition := range partitionList {
//        // 针对每个分区创建一个对应的分区消费者
//        // 消费 web_log 的 partition 分区，从 offset 开始             offset的意义？？
//        pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
//        if err != nil {
//            log.Println(partition, err)
//            return
//        }
//        defer pc.AsyncClose()
//        // 异步从每个分区消费信息
//        go func(sarama.PartitionConsumer) {
//            for msg := range pc.Messages() {
//                fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
//            }
//        }(pc)
//    }
//    time.Sleep(10*time.Second)
//}













































