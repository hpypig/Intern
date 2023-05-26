package main
// ref：https://www.lixueduan.com/posts/kafka/07-partition/
import "github.com/Shopify/sarama"

//type Config struct {
//    Producer struct {
//        Partitioner PartionerConstructor
//    }
//}

//type PartitionerConstructor func(topic string) Partitioner

//type Partitioner interface {
//    Partition(message *ProducerMessage, numPartitions int32) (int32, error)
//    RequiresConsistency() bool
//}

type myPartitioner struct {
    partition int32
}
// Partition 返回的是分区的位置或者索引，并不是具体的分区号。
// 比如有十个分区[0,1，2,3...9] 这里返回 0 表示取数组中的第0个位置的分区。
// numPartitions 是打算使用的分区数（指定一致性时，不可用的分区也要往里发）
func (p *myPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
    if p.partition >= numPartitions {
        p.partition = 0
    }
    ret := p.partition
    p.partition++
    return ret, nil
}
func (p *myPartitioner) RequiresConsistency() bool {
    return false // 不要一致性，即分区不可用则不往里发？
}
func NewMyPartitioner(topic string) sarama.Partitioner {
    return &myPartitioner{}
}




//type tt struct {
//    stu struct {  // 这是什么写法
//        a int
//    }
//}

func t() {
}
