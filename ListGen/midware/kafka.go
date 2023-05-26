package midware

import (
    "ListGen/entities"
    "encoding/json"
    "github.com/Shopify/sarama"
    "log"
)

const (
	Topic = "UpdatedData"
    OldOffset = sarama.OffsetOldest
    NewOffset = sarama.OffsetNewest
)

func GetUpdatedData(ch chan entities.UpdatedDataResponse) {
    config := sarama.NewConfig()
    consumer, err := sarama.NewConsumer([]string{"120.79.29.70:9092"}, config)
    if err != nil {
        log.Fatal("GetUpdatedData-NewConsumer err: ", err)
    }
    defer consumer.Close()
    //partitionConsumer, err := consumer.ConsumePartition(Topic, 0, NewOffset) // 增量更新时用
    partitionConsumer, err := consumer.ConsumePartition(Topic, 0, OldOffset) // 初次更新数据时用
    if err != nil {
        log.Fatal("ConsumePartition err: ", err)
    }
    defer partitionConsumer.Close()
    for message := range partitionConsumer.Messages() {
        log.Printf("[Consumer] partitionid: %d; offset:%d", message.Partition, message.Offset)
        var updatedData entities.UpdatedDataResponse
        err = json.Unmarshal(message.Value, &updatedData)
        if err != nil {
            log.Println("GetUpdatedData-unmarshal err: ", err)
            return
        }
        log.Printf("GetUpdatedData-updatedData id: %v; title:%v\n", updatedData.Id, updatedData.Data.Title )
        ch <- updatedData
    }

}
