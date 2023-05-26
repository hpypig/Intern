package main

import (
    "fmt"
    "gopkg.in/redis.v3"
    "log"
)
var Client *redis.Client
/*
SortedSet
1）添加
2）修改
3）移除 30 名以外
4）读取整个集合


*/
func main() {
    Connect()
    //SortedSet()
    //SortedSetAdd()
    //SortedSetGet()
    //ZRemove("zset1","id1")
    example()
}


func Connect() {
    Client = redis.NewClient(&redis.Options{
        Addr:     "120.79.29.70:6379",
        Password: "710710", // no password set
        DB:       0,  // use default DB
    })

    pong, err := Client.Ping().Result()
    fmt.Println(pong, err)
    // Output: PONG <nil>
}

func example() {
    // key value expiration
    statusCmd := Client.Set("akey0","testvalue",0)
    err := statusCmd.Err()
    if err != nil {
        log.Println("err: ", err)
        return
    }

    any1, err := statusCmd.Result()
    fmt.Println(any1,err) // OK <nil>

    str := statusCmd.String()
    any2 := statusCmd.Val()
    fmt.Println(str) // SET testkey testvalue: OK
    fmt.Println(any2) // OK
}

// SortedSetAdd 插入
func SortedSetAdd() {
    for i:=0; i<=10; i++ {
        err := Client.ZAdd("testkey2",redis.Z{float64(10+i),"id"+fmt.Sprintf("%d",i)}).Err()
        if err != nil {
            log.Println("zadd err: ", err)
        }
    }
}
// SortedSetGet 查询
func SortedSetGet() {
    key := "testkey2"
    var start,stop int64
    start, stop = 0, -1
    //stringSliceCmd := Client.ZRange(key, start, stop)
    //if err := stringSliceCmd.Err(); err != nil {
    //    log.Println("stringSliceCmd err: ",err)
    //    return
    //}
    //str := stringSliceCmd.String()
    //fmt.Println(str) // ZRANGE testkey2 0 -1: [id0 id1 id2 id3 id4 id5 id6 id7 id8 id9 id10]
    //strs := stringSliceCmd.Val()
    //fmt.Println(strs) // [id0 id1 id2 id3 id4 id5 id6 id7 id8 id9 id10]
    //strs2, err := stringSliceCmd.Result() // [id0 id1 id2 id3 id4 id5 id6 id7 id8 id9 id10] <nil>
    //fmt.Println(strs2, err)
    // ----------------------
    stringSliceCmd2 := Client.ZRangeWithScores(key, start, stop)
    if err := stringSliceCmd2.Err(); err != nil {
        log.Println("stringSliceCmd err: ",err)
        return
    }
    str := stringSliceCmd2.String()
    zs := stringSliceCmd2.Val()
    zs2,err := stringSliceCmd2.Result()
    fmt.Println(str)
    fmt.Println(zs)
    fmt.Println(zs2,err)
}

func ZRemove(key string, member string) {
    intCmd := Client.ZRem(key, member)
    num := intCmd.Val()
    err := intCmd.Err()
    fmt.Println(num, err) // 1 <nil>

    num, err = intCmd.Result()
    fmt.Println(num, err) // 1 <nil>

    str := intCmd.String()
    fmt.Println(str) // ZREM zset1 id1: 1
}
func SRem(key string, member string) {
    //intCmd := Client.SRem(key, member)
}




