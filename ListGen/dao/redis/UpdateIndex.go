package redis

import (
    "errors"
    "fmt"
    "gopkg.in/redis.v3"
    "log"
    "time"
)

// 向 stock 索引 添加 score id

func PostIDZSet(key string, publishTime int64, id string) (int64, error){ // 这个 int64 是什么
    //key = strings.ToUpper(key)
    fmt.Printf("key: %s publishTime: %d id: %s\n",key,publishTime,id)
    return client.ZAdd(key, redis.Z{float64(publishTime), id}).Result()
}

func SAddIDToKeys(id string, key string) (res int64, err error) {
    err = errors.New("SAddIDToKeys err")

    res, err = client.SAdd(id, key).Result()
    if err != nil {
        log.Println("id to key err: ", err)
        return
    }
    // 返回的布尔值还不知道是用来干嘛的
    _, err = client.Expire(id, time.Hour*24*30).Result()
    if err != nil {
        log.Println("id to key err: ", err)
        return
    }
    return res, nil
}

func SetLastTime(lastTime int64) (res string, err error ) {
    return client.Set("last_time", lastTime,time.Hour*24*30).Result()
}

//------------

func PostIDTitleMediaZSet(key string, publishTime int64, id string, title string, media string) (int64, error){ // 这个 int64 是什么
    //key = strings.ToUpper(key)
    fmt.Printf("key: %s publishTime: %d id: %s\n",key,publishTime,id)
    return client.ZAdd(key, redis.Z{float64(publishTime), id+"_"+title+"_"+media}).Result()
}









