package redis

func SMembersByID(key string, members *[]string) (err error) {
    *members, err = client.SMembers(key).Result()
    return err
}

func GetLastTime() (lastTime int64, err error) {
    // ttl key 查看剩余过期时间
    //client.Get("last_time").Result()
    // 有bug，不知道为什么有时查出来是nil，明明没过期；这种情况下更新了redis后是36，37没有更新
    return client.Get("last_time").Int64()
}

func GetIdsByRevRange(key string, start int64, end int64) ([]string, error) {
    //fmt.Printf("GetIdsByKey key: %s\n", key)
    return client.ZRevRange(key, start, end).Result()
}
