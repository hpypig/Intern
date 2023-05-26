package mongo

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "log"
)

var session *mgo.Session

func Init() {
    var err error
    URL := "mongodb://120.79.29.70:30001,120.79.29.70:30002,120.79.29.70:30003/?replicaSet=my-mongo-set"
    session, err = mgo.Dial(URL)
    if err != nil {
        log.Println("dial err: ", err)
        return
    }
    err = session.Ping()
    if err != nil {
        log.Println("ping err: ",err)
        return
    }
    fmt.Println("connect success")
}
func Close() {
    session.Close()
}
