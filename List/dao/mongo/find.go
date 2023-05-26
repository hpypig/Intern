package mongo

import (
    "fmt"
    "github.com/hpypig/Intern/entities"
    "gopkg.in/mgo.v2/bson"
    "log"
    "time"
)


func FindContentByObjectID(_id bson.ObjectId, content *entities.NewsContent) error { // 像这种就应该单独测试，但需要connect，就比较麻烦
    sess := session.Copy()
    defer sess.Close()
    query := bson.M{"_id":_id}
    return sess.DB("gf").C("NewsContent").Find(query).One(content)

}




func IterTest() (err error) {
    sess := session.Copy()
    defer sess.Close()

    collection := sess.DB("local").C("oplog.rs")


    //query := bson.M{"ns":"gf.NewsContent", "ts":bson.M{"$gte":bson.MongoTimestamp(7175036832709607425)}}
    query := bson.M{"ns":"gf.NewsContent", "ts":bson.M{"$gte":bson.MongoTimestamp(1670746737<<32)}}
    //var oplog []entities.Oplog
    //collection.Find(query).Limit(2).All(&oplog)
    //fmt.Printf("%+v\n",oplog)

    // -1 next一直等，只要cursor合法且session没有结束
    // 0 在请求的集合访问完时就 timeout，不会等
    // 可设置时间，访问完后等多久
    var oplog entities.Oplog
    iter := collection.Find(query).Limit(20).Tail(2*time.Second) //为什么不设置成阻塞式的
    for {
        for iter.Next(&oplog) {
            fmt.Printf("%+v\n",oplog)
            //time.Sleep(5*time.Second)
        }
        // iter结构体内部设置了err字段，并且预先定义了常量error，会在iter出错时赋给该字段
        // 操作过程加了锁，为什么要加锁？有没有什么快捷办法能找到各个临界区？搜索？
        // 下面两行的作用就是判断 Next 结束的原因，做相应处理
        if err := iter.Err(); err != nil {return iter.Close()}
        if iter.Timeout() {continue}
    }


    //var oplog entities.Oplog
    //query := bson.M{"ns":"gf.NewsContent","ts":bson.M{"$gt":7175036832709607425}}
    ////query := bson.M{"ns":"gf.NewsContent"}
    //
    ////iter := collection.Find(nil).Sort("$natural").Tail(5 * time.Second)
    //iter := collection.Find(query).Sort("$natural").Tail(3*time.Second)
    //for {
    //   for iter.Next(&oplog) {
    //       fmt.Printf("%+v\n",oplog)
    //   }
    //   if iter.Err() != nil {
    //       return iter.Close()
    //   }
    //   if iter.Timeout() {
    //       log.Println("timeout")
    //       return
    //   }
    //
    //   //query := collection.Find(bson.M{"_id": bson.M{"$gt": lastId}})
    //   //iter = query.Sort("$natural").Tail(5 * time.Second)
    //}
    //iter.Close()
    return nil
}


func GetBriefInfoById(id string, data *entities.BriefInfo) (err error){
    sess := session.Copy()
    defer sess.Close()
    query := bson.M{"id":id}
    err = sess.DB("gf").C("NewsContent").Find(query).One(data)
    return
}

//func GetBriefInfoByIds(ids []string) []entities.BriefInfo {
//    data := make([]entities.BriefInfo, len(ids))
//    // 根据 id 列表取出标题等数据
//    for i,id := range ids {
//        err := GetBriefInfoById(id, &data[i])
//        if err != nil {
//            //log.Printf("id: %s; GetBriefInfoById err: %v\n", id, err)
//            fmt.Printf("id: %s; GetBriefInfoById err: %v\n", id, err)
//        }
//        //fmt.Printf("%+v\n", data[i])
//    }
//    //fmt.Println("find-GetBriefInfoByIds out")
//    return data
//}

// 因为性能差，估计是因为每个id都copy的原因？

func GetBriefInfoByIds(ids []string) []entities.BriefInfo {
    sess := session.Copy()
    defer sess.Close()
    data := make([]entities.BriefInfo, len(ids))
    // 根据 id 列表取出标题等数据
    //t1 := time.Now()
    for i,id := range ids {
        //t2 := time.Now()
        query := bson.M{"id":id}
        //fields := bson.M{"id":1, "title":1, "media":1, "stocks":1, "publishTime":1}
        err := sess.DB("gf").C("NewsContent").Find(query).One(&data[i])
        //err := sess.DB("gf").C("NewsContent").Find(query).Select(fields).One(&data[i])

        if err != nil {
            //log.Printf("id: %s; GetBriefInfoById err: %v\n", id, err)
            fmt.Printf("id: %s; GetBriefInfoById err: %v\n", id, err)
        }
        //fmt.Printf("%+v\n", data[i])
        //fmt.Println("GetBriefInfoByIds-single 时间消耗: ",time.Since(t2))
    }
    //fmt.Println("GetBriefInfoByIds 时间消耗: ",time.Since(t1))
    //fmt.Println("find-GetBriefInfoByIds out")
    return data
}

func GetBriefInfoByIds2(ids []string) []entities.BriefInfo {
    sess := session.Copy()

    defer sess.Close()
    var data []entities.BriefInfo
    query := bson.M{"id":bson.M{"$in":ids}}
    fields := bson.M{"id":1, "title":1, "media":1, "stocks":1, "publishTime":1}
    //err := sess.DB("gf").C("NewsContent").Find(query).All(&data)
    err := sess.DB("gf").C("NewsContent").Find(query).Select(fields).All(&data)

    if err != nil {
        fmt.Printf("ids: %v; GetBriefInfoById err: %v\n", ids, err)
    }
    return data
}


func GetNewsContentById(id string, data *entities.NewsContent){
    sess := session.Copy()
    defer sess.Close()
    query := bson.M{"id":id}
    err := sess.DB("gf").C("NewsContent").Find(query).One(data)
    if err != nil {
        log.Println("GetNewsContentById err:", err)
        return
    }
}
