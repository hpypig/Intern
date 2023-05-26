package redis

import (
    "fmt"
    "gopkg.in/redis.v3"
    "log"
)

const (
    StockListKeyPrefix = "info_stock_list"
    ImportantNewsListKey = "info_important_list"
    ColumnListKeyPrefix = "info_column_list"
    IndexNameListPrefix = "index_name_list"
)


var client *redis.Client
func Init() {
    client = redis.NewClient(&redis.Options{
        Addr:     "120.79.29.70:6379",
        //Addr:     "127.0.0.1:6379",
        Password: "710710", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    if err != nil {
        log.Println("ping err: ", err)
    }
    fmt.Println("pong: ",pong)
    // Output: PONG <nil>
}

func Close() {
    client.Close()
}
