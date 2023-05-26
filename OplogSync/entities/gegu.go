package entities
type OptionalDetail struct {  // 自选部分 0123 的细节 node
    Id string `bson:"id"`
    Title string `bson:"title"`
    PublishTime int64 `bson:"publishTime"`
    Media string       `bson:"media"`
    Stocks []Stock     `bson:"stocks"`
    Columns []string   `bson:"columns"`
    TxtType int `bson:"txtType"`// ???
    Industries []Stock `bson:"industries"`
}
// test: id = "638d9e50e4f8920007314bac"
