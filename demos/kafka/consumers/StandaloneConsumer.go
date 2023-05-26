package main

import (
    "fmt"
    "github.com/Shopify/sarama"
    "log"
    "sync"
    "time"
)

/*
1）生成 config
2）建立连接，得到 consumer
3）defer close
4）消费，指定 topic partition offset     指定分区有什么实际意义？难道我每次投消息都要记录哪个分区有什么，然后消费时再指定？
*/

// SinglePartition 单分区消费
func SinglePartition(topic string) {
    config := sarama.NewConfig()
    consumer, err := sarama.NewConsumer([]string{"120.79.29.70:9092"}, config)
    if err != nil {
        log.Fatal("NewConsumer err: ", err)
    }
    defer consumer.Close()
    // 参数1 指定消费哪个 topic
    // 参数2 分区 这里默认消费 0 号分区 kafka 中有分区的概念，类似于ES和MongoDB中的sharding，MySQL中的分表这种
    // 参数3 offset 从哪儿开始消费起走，正常情况下每次消费完都会将这次的offset提交到kafka，然后下次可以接着消费，
    // 这里demo就从最新的开始消费，即该 consumer 启动之前产生的消息都无法被消费
    // 如果改为 sarama.OffsetOldest 则会从最旧的消息开始消费，即每次重启 consumer 都会把该 topic 下的所有消息消费一次
    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatal("ConsumePartition err: ", err)
    }
    defer partitionConsumer.Close()
    //会一直阻塞在这里 ？？？？？？？？？？？？？？？？？？？？？？？？？为啥会阻塞？？？？
    for message := range partitionConsumer.Messages() {
       log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, string(message.Value))
    }
    // 一般不会超时退出，因为要等生产者继续发
    //t := time.NewTimer(3*time.Second)
    //for {
    //    t.Reset(3*time.Second)
    //    select {
    //    case message := <-partitionConsumer.Messages():
    //        log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, string(message.Value))
    //    case <-t.C:
    //        fmt.Println("3s")
    //        return
    //    }
    //}

}

/*
1）config
2）consumer
3）defer close
4）根据 topic 得到分区 partitions []int32
5）根据分区数 wg.Add
6）对每个分区进行消费
    使用 consumer 生成 partitionConsumer
    使用 partitionConsumer 获取 <-chan *ConsumerMessage
*/

// Partitions 多分区消费
func Partitions(topic string) {
    config := sarama.NewConfig()
    consumer, err := sarama.NewConsumer([]string{"120.79.29.70:9092"}, config)
    if err != nil {
        log.Fatal("NewConsumer err: ", err)
    }
    defer consumer.Close()
    partitions, err := consumer.Partitions(topic)
    if err != nil {
        log.Fatal("Partitions err: ", err)
    }
    var wg sync.WaitGroup
    wg.Add(len(partitions))
    for _, partitionID := range partitions {
        go consumeByPartition(consumer, topic, partitionID, &wg)
    }
    wg.Wait()
}
func consumeByPartition(consumer sarama.Consumer, topic string, partitionID int32, wg *sync.WaitGroup) {
    defer wg.Done()
    partitionConsumer, err := consumer.ConsumePartition(topic, partitionID, sarama.OffsetOldest) // 从头消费
    if err != nil {
        log.Fatal("ConsumePartition err: ", err)
    }
    defer partitionConsumer.Close()
    // 为啥这里总是会阻塞？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？？只能超时控制？
    for message := range partitionConsumer.Messages() {
        log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, string(message.Value))

    }
}

/*
Kafka 和其他 MQ 最大的区别在于 Kafka 中的消息在消费后不会被删除，而是会一直保留，直到过期。
为了防止每次重启消费者都从第 1 条消息开始消费，我们需要在消费消息后将 offset 提交给 Kafka。
这样重启后就可以接着上次的 Offset 继续消费了。
*/

/*
1）配置 offset 自动提交、和提交时间间隔
2）连接，生成 Client

3）创建偏移量管理器
4）创建对应分区的偏移量管理器
5）记录偏移量
6）提交偏移量

*/

func OffsetManager(topic string) {
    config := sarama.NewConfig()
    config.Consumer.Offsets.AutoCommit.Enable = true
    config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
    client, err := sarama.NewClient([]string{"120.79.29.70:9092"}, config)
    if err != nil {
        log.Fatal("NewClient err: ", err)
    }
    defer client.Close()
    // offsetManager 用于管理每个 consumerGroup的 offset
    // 根据 groupID 来区分不同的 consumer，注意: 每次提交的 offset 信息也是和 groupID 关联的
    offsetManager, err := sarama.NewOffsetManagerFromClient("0", client)
    if err != nil {
        log.Println("NewOffsetManagerFromClient err:", err)
    }
    defer offsetManager.Close()
    // 每个分区的 offset 也是分别管理的，demo 这里使用 0 分区，因为该 topic 只有 1 个分区
    partitionOffsetManager, err := offsetManager.ManagePartition(topic, 0)
    if err != nil {
        log.Println("ManagePartition err:", err)
    }
    defer partitionOffsetManager.Close()
    // defer 在程序结束后在 commit 一次，防止自动提交间隔之间的信息被丢掉
    defer offsetManager.Commit()
    consumer, err := sarama.NewConsumerFromClient(client)  // 感觉这是某种设计模式，可以深入研究
    if err != nil {
        log.Println("NewConsumerFromClient err:", err)
    }
    // 根据 kafka 中记录的上次消费的 offset 开始+1的位置接着消费
    nextOffset, _ := partitionOffsetManager.NextOffset()
    fmt.Println("nextOffset:", nextOffset)
    partitionConsumer, err := consumer.ConsumePartition(topic, 0,nextOffset)
    if err != nil {
        log.Println("ConsumePartition err:", err)
    }
    defer partitionConsumer.Close()
    for message := range partitionConsumer.Messages() {
        value := string(message.Value)
        log.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, value)
        // 每次消费后都更新一次 offset,这里更新的只是程序内存中的值，需要 commit 之后才能提交到 kafka
        partitionOffsetManager.MarkOffset(message.Offset+1, "modified metadata") // MarkOffset 更新最后消费的 offset
    }
}
