package redis

import (
    "fmt"
    "testing"
)

func TestSetLastTime(t *testing.T) {
    Init()
    res, err := SetLastTime( 1670746737)
    if err!=nil {
       fmt.Println(err)
    }
    fmt.Println(res)

    // 当键值为空时返回 err 和 0
    //res2, err := GetLastTime()
    //if err != nil {
    //    fmt.Println("GetLastTime err: ", err) // GetLastTime err:  redis: nil
    //}
    //fmt.Println(res2) // 0

}
