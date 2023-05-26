package mongo

import (
    "fmt"
    "ListGen/entities"

    //"github.com/hpypig/Intern/entities"
    //"gopkg.in/mgo.v2/bson"
    "testing"
)

//func TestParseOplogAfter(t *testing.T) {
//    Init()
//    var content entities.NewsContent
//    err := FindContentByObjectID(bson.ObjectIdHex("63959271b63543b3ed9c5d26"), &content)
//    if err != nil {
//        fmt.Println(err)
//        return
//    }
//    fmt.Printf("%+v\n",content)
//    //ch := make(chan entities.Oplog, 10)
//    //ParseOplogAfter(0,ch)
//}
func TestFindContentByObjectID(t *testing.T) {
    Init()
    id := "63947a1f86ded60007490b94"
    var info entities.BriefInfo
    GetBriefInfoById(id, &info)
    fmt.Printf("res: %+v\n", info)
}
