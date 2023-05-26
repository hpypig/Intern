package redis
func ZRemove(key string, members ...string) (int64, error) {
    return client.ZRem(key, members...).Result()
}

func SRemove(key string, member ...string) (int64, error) {
    return client.SRem(key,member...).Result()
}

func Scan() {
    //client.Scan()
}
