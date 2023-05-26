package main
/*
ref:
   http://learningprogramming.net/golang/golang-and-mongodb/group-by-and-sort-in-golang-and-mongodb/
   https://blog.kelu.org/tech/2020/11/01/golang-mgo-mongodb.html
   https://pkg.go.dev/gopkg.in/mgo.v2


*/
import (
    "encoding/json"
    "fmt"
    "github.com/hpypig/Intern/demos/mgov2/entities"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "io/ioutil"
    "log"
    "net/http"
)

const (
    URL = "mongodb://120.79.29.70:30001,120.79.29.70:30002,120.79.29.70:30003/?replicaSet=my-mongo-set"
    //DBName = "gf"
)


var DB *mgo.Database
var Session *mgo.Session  // 通过复制不用开多个连接？？



func main() {
    //url := "mongodb://ww:123111@120.79.29.70:27017/jinyun?authMechanism=Scram" // &connect=direct
    //url := "mongodb://120.79.29.70:27017?connect=direct"
    // 不指定复制集，采用直连方式，主从切换后无法写入（只能在开发环境用直连）
    DBName := "gf"
    //var err error
    // 连接，获取数据库
    err := Connect(URL)
    if err != nil {
        return
    }
    DB = Session.DB(DBName)
    //FindOne("ImportantNewsNodes","638f2ddee4f89200073157fa")
    //FindContentById("NewsContent","638f2ddee4f89200073157fa")

    //FindOptionalDetailById("StockNodes", "638d9e50e4f8920007314bac") // 查自选股相关文章的基本信息
    //FindContentById("NewsContent","638d9e50e4f8920007314bac")

    //ids := []string{"638d9e50e4f8920007314bac","638f2ddee4f89200073157fa","638f0c2be4f8920007315685"}
    //FindManyByIds("NewsContent",ids)

    //DataPrepare("ImportantNewsNodes")

    //Insert() //测试用
    // 通过DataPrepare获取数据后插入测试集合
    InsertNewsContent("ImportantNewsNodes")
    InsertNewsContent("StockNodes")
    InsertNewsContent("ColumnNodes")


    Session.Close()


}
func Connect(url string) (err error) {
    Session, err = mgo.Dial(url)
    if err != nil {
        log.Println("dial err: ", err)
        return
    }
    err = Session.Ping()
    if err != nil {
        log.Println("ping err: ",err)
        return
    }
    return nil
}


func FindOne (collection string, id string) (err error) {
    //session := Session.Copy()
    c := DB.C(collection)
    detail := entities.ImportantNewsDetail{}
    err = c.Find(bson.M{"id":id}).One(&detail)
    if err != nil {
        log.Println("FindOne err:",err)
        return
    }
    fmt.Printf("%+v\n", detail)



    return nil
}
func FindContentById(collection string, id string) (content entities.NewsContent,err error) {
    c := DB.C(collection)
    content = entities.NewsContent{}
    err = c.Find(bson.M{"id":id}).One(&content)
    if err != nil {
        log.Println("FindOne err:",err)
        return
    }
    //fmt.Printf("%+v\n", content)
    return content,nil
}

// "638d9e50e4f8920007314bac"

// FindOptionalDetailById 寻找自选股文章的基本信息
func FindOptionalDetailById(collection string, id string) (err error) {
    //session := Session.Copy()
    c := DB.C(collection)
    var detail entities.OptionalDetail
    err = c.Find(bson.M{"id":id}).One(&detail)
    if err != nil {
        log.Println("FindOptionalDetailById err:",err)
        return
    }
    fmt.Printf("%+v\n", detail)
    return nil
}

// 如何根据 id 集合返回一堆内容，并且拼接成一个新的结构体，然后传给前端？

func FindManyByIds(collection string, ids []string) {
    c := DB.C(collection)
    var articles []entities.NewsContent
    c.Find(bson.M{"id":bson.M{"$in":ids}}).All(&articles)
    for _,v := range articles {
        fmt.Printf("%+v\n", v)
    }
}

//----------------插入内容数据

type ID struct {
    Id string `bson:"id"`
}
// DataPrepare 查询各索引 id，然后请求 origin 接口数据，然后插入数据库。
func DataPrepare(collection string) ([]entities.PrepareData,error) {
    c := DB.C(collection)
    var ids []ID
    // 从mongodb取出节点 id ;  只取 id 字段
    c.Find(bson.M{}).Select(bson.M{"id":1}).All(&ids)
    fmt.Println(len(ids),ids)
    // 638f2ddee4f89200073157fa 该 id 对应的内容不再插入
    // 逐 id 调用 origin 接口，将获取的数据插入db
    url := "https://info.gf.com.cn/api/1.0.0/read/article/"
    var res []entities.PrepareData
    for _,v := range ids {
       //fmt.Println(url+v.Id)
       resp, err := http.Get(url+v.Id)
       if err != nil {
           fmt.Println("DataPrepare http get err: ", err)
           return nil, err
       }

       body, err := ioutil.ReadAll(resp.Body)
       //fmt.Println(len(body))
       if err != nil {
           log.Println("body read err: ", err)
           return nil, err
       }
       var prepareData entities.PrepareData
       json.Unmarshal(body, &prepareData)
       //fmt.Printf("%+v\n",prepareData.Data)
       res = append(res, prepareData)
    }
    fmt.Println("count:", len(res))
    return res, nil
}
func SwitchDB(db string) {
    DB = Session.DB(db)
}
func Insert() {
    SwitchDB("test")
    contentCollection := DB.C("content")
    contentCollection.Insert(&entities.PrepareData{Data:entities.NewsContent{Source:12}})
}

func InsertNewsContent(collection string) {
    // 在gf，获取 collection 索引，得到origin内容
    data, err := DataPrepare(collection)
    if err != nil {
        log.Println("****")
        return
    }

    //c := DB.C(collection)
    // 切换数据库到test，插入测试数据集 content
    //SwitchDB("test")

    contentCollection := DB.C("NewsContent")
    for i,_ := range data {
        err = contentCollection.Insert(&data[i].Data)
        if err != nil {
            log.Println("insert err: ",err)
            return
        }
    }


}


