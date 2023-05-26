package entities

import (
    "gopkg.in/mgo.v2/bson"
    "time"
)

type User struct {
    Id_       bson.ObjectId `bson:"_id"`
    Name      string        `bson:"name"`
    Age       int           `bson:"age"`
    JonedAt   time.Time     `bson:"joned_at"`
    Interests []string      `bson:"interests"`
}

type BriefInfo struct {  // Node  这个里面的信息就够了，似乎不需要具体内容？那update呢？
    Id string `bson:"id" json:"id"`
    Title string `bson:"title" json:"title"`
    Media string      `bson:"media" json:"media"`
    Stocks []Stock    `bson:"stocks" json:"stocks"`
    PublishTime int64 `bson:"publishTime" json:"publishTime"` // 时间戳用int32位表示接收后得到的是负值，uint32也不对
}
type Stock struct {
    Market string `bson:"market"`
    Code string `bson:"code"`
    Name string `bson:"name"`
    Type int `bson:"type"` // 类型？？要闻里的内容有 tag 无 type 字段
    Tag StockTag `bson:"tag"`
}

// 小数字设置什么类型？？

type StockTag struct {
    Weight float64 `bson:"weight"`
    Emotion string `bson:"emotion"`
    EmotionWeight int64 `bson:"emotionWeight"`
}

// 各个id返回的 NewsContent 字段其实不太一样。比如：有的有stocks，有的没stocks

type NewsContent struct { // 要闻的内容结构和自选股文章内容结构应该是相同的
    Source int64 `bson:"source" json:"source"`
    Id string `bson:"id" json:"id"`
    Title string `bson:"title" json:"title"`
    Subtitle string `bson:"subtitle" json:"subtitle"`
    Media string `bson:"media" json:"media"`
    Content string `bson:"content" json:"content"`
    //Privilege int
    Status int `bson:"status" json:"status"`
    CreateTime int64 `bson:"createTime" json:"createTime"`
    UpdateTime int64 `bson:"updateTime" json:"updateTime"`
    PublishTime int64 `bson:"publishTime" json:"publishTime"`
    Categories []string `bson:"categories,omitempty" json:"categories"`
    Columns []string  `bson:"columns,omitempty" json:"columns"`
    Stocks []Stock    `bson:"stocks,omitempty" json:"stocks"`  // 有的内容没有股票
    SourceName string `bson:"sourceName" json:"sourceName"`
    ColumnsObj []ColumnsObjection `bson:"columnsObj" json:"columnsObj"`
    TxtType int64  `bson:"txtType" json:"txtType"`// ?????用什么类型  0 1 2 3
    Links []Link `bson:"links,omitempty" json:"links"`
    //tts
    MaskTitle string `bson:"maskTitle" json:"maskTitle"`
}
type ColumnsObjection struct {
    Id string `bson:"id" json:"id"`
    Title string `bson:"title" json:"title"`
}
type Link struct {
    Word string `bson:"word" json:"word"`
    Type int `bson:"type" json:"type"`
    Target string `bson:"target" json:"target"`
}
type ImportantNewsReturn struct {

}


//--------------------- 用于提取 origin 接口的数据内容

type PrepareData struct {
    //ErrorCode int `json:"errorCode"`
    //CostTime int `json:"costTime"`
    Data NewsContent `json:"data"`  // 犯的错：误把结构体定位切片，导致读出为空
}

//------------ 提取oplog----问题：现在是 NewsContent 写死的，以后可能要用 interface

//type Oplog struct {
//    Ts bson.MongoTimestamp `bson:"ts"`
//    //Ts int64 `bson:"ts"`  // 接收可以用 int64，query 不可以
//    O NewsContent `bson:"o"`
//}

type Oplog struct {
    Ts int64 `bson:"ts"`  // 接收用 int64，发送时要用 MongoTimeStamp
    Ns string `bson:"ns"`   // gf.NewsContent
    //Ts int64 `bson:"ts"`  // 接收可以用 int64，query 不可以
    Op string `bson:"op"`   // u、i
    O ForID `bson:"o"`  // i 时 _id 在这里面
    O2 ForID `bson:"o2"` // u 时 _id 在这里面

}
type ForID struct {
    Id bson.ObjectId `bson:"_id"`
}

type UpdatedDataResponse struct {
    Op string `json:"op"`
    Id string `json:"id"`
    Data NewsContent `json:"data"`
}
