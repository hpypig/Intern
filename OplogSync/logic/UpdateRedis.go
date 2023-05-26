package logic

import (
    "fmt"
    "log"
    "oplogsync/dao/redis"
    "oplogsync/entities"
)


/*
   根据 id 取出数据，更新索引
       对应 key 处添加 id
       更新 id 关联的 keys
       根据删除的 keys 在对应 key 删除该 id(插入操作才需要删除)

    记录：做了拆分操作，把接收 oplog 并查 数据拆分出去单独写
*/

// UpdateRedis 更新 Redis 索引
func UpdateRedis(updatedDataChan <-chan entities.UpdatedDataResponse) { // logic 层
    for {
        updatedData := <-updatedDataChan
        content := updatedData.Data
        // 1）个股：资讯类型、市场、代码  发布时间  id
        keyMap := map[string]bool{}
        // 添加 id
        for _,stock := range content.Stocks {
            // 按资讯类型添加 id
            StockListKey := redis.StockListKeyPrefix + fmt.Sprintf("_%d_%s_%s",content.TxtType, stock.Market, stock.Code)
            res, err := redis.PostIDZSet(StockListKey, content.PublishTime, content.Id)  // 向 stock 索引 添加 score id
            if err != nil {
                log.Println("StockList err: ", err)
                return
            }
            fmt.Printf("股票id列表添加结果: %d\n", res)
            keyMap[StockListKey] = true // id 当前 keys

            // 对全集资讯列表添加id
            StockListKey2 := redis.StockListKeyPrefix + fmt.Sprintf("_%d_%s_%s",0, stock.Market, stock.Code)
            res2, err := redis.PostIDZSet(StockListKey2, content.PublishTime, content.Id)  // 向 stock 索引 添加 score id
            fmt.Printf("0型列表添加结果: %d\n", res2)
            keyMap[StockListKey2] = true

            // 维护 id 对应股票名（索引名）-- 1、2、3 类
            res, err = redis.SAddIDToKeys(redis.IndexNameListPrefix + "_" + content.Id, StockListKey)
            if err != nil {
                log.Println("SAddIDToKeys err: ", err)
                return
            }
            fmt.Printf("SAddIDToKeys res: %d\n", res)
            // 0 类（代表全部）
            res2, err = redis.SAddIDToKeys(redis.IndexNameListPrefix + "_" + content.Id, StockListKey2)
            if err != nil {
                log.Println("SAddIDToKeys2 err: ", err)
                return
            }
            fmt.Printf("SAddIDToKeys2 res: %d\n", res2)
        }

        // 如果是更新操作，要多做一步删除，删除旧 keys 中和 id 解关联的 key，删除 id 列表中，解关联的 id
        // 不管 i、u 我上面都会插入redis一次，用map记录插入的索引名（股票）
        // 遍历 redis 存的 idToName，检查哪些不存在于新 map，哪些就要被删除
        if updatedData.Op == "u" {
            // 删除已不在当前keys列表的key：得到keys set，把不在 map 内的删除
            var keys []string
            err := redis.SMembersByID(content.Id, &keys)
            if err != nil {
                log.Println("SMembersByID err: ",err)
                return
            }
            fmt.Printf("id: %s 关联的旧股票列表 keys:%v \n", content.Id, keys)

            // 把该 id 从某些索引中删除(即上面解除关联的键)
            for _,key := range keys {
                if !keyMap[key] { // 如果新 map 没有此 key
                    // 在 stock 索引删除此 id
                    res, err := redis.ZRemove(key,content.Id)
                    if err != nil {
                        log.Printf("ZRemove err:%v, key:%s, id:%s\n", err, key, content.Id)
                        return
                    }
                    fmt.Printf("ZRemove res:%v, key:%s, id:%s\n", res, key, content.Id)

                    // 在历史 stock 集合删除此 stock（key）
                    res, err = redis.SRemove(redis.IndexNameListPrefix + "_" + content.Id, key)
                    if err != nil {
                        log.Printf("ZRemove err:%v, key:%s, id:%s\n", err, key, content.Id)
                        return
                    }
                    fmt.Printf("ZRemove res:%v, key:%s, id:%s\n", res, key, content.Id)
                }
            }
        }
        // zrem info_important_list 6394603486ded60007490by8 测试用
        // 2）要闻：id
        if content.TxtType == 1 {
            res, err := redis.PostIDZSet(redis.ImportantNewsListKey, content.PublishTime, content.Id)
            if err != nil {
                log.Println("ImportantNews err: ", err)
                return
            }
            fmt.Printf("ImportantNews res: %d\n", res)
            // 要删除 30 条以后的（倒序）
            // ZREMBYRANK key 0 -31


        }
        // 3）栏目：栏目id
        for _,column := range content.Columns {
            ColumnListKey := redis.ColumnListKeyPrefix + "_" + column
            res, err := redis.PostIDZSet(ColumnListKey, content.PublishTime, content.Id)
            if err != nil {
                log.Println("ColumnList err: ", err)
                return
            }
            fmt.Printf("ColumnList res: %d\n", res)
        }

    }
}


