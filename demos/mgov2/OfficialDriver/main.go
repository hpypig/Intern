package main
import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

// ref: https://www.liwenzhou.com/posts/Go/mongodb/#autoid-0-4-4

var Client *mongo.Client
//var DB *mongo.Database

func Connect() {
    //clientOptions := options.Client().ApplyURI("mongodb://120.79.29.70:27017")

    credential := options.Credential{
       Username:      "ww",
       Password:      "123111",
    }
    clientOptions := options.Client().ApplyURI("mongodb://120.79.29.70:27017").SetAuth(credential)
    //Client, err := mongo.Connect(
    //    context.TODO(),
    //    clientOptions,
    //    options.Client().SetAuth(awsCredential))

    // 连接到MongoDB
    c, err := mongo.Connect(context.TODO(), clientOptions)
    Client = c
    if err != nil {
        log.Fatal(err)
    }

    // 检查连接
    err = Client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")




}
func Close() {
    err := Client.Disconnect(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connection to MongoDB closed.")
}
func getDB(dbName string) (DB *mongo.Database) {
    return Client.Database(dbName)
}


func insertOne(db *mongo.Database) {
    collection := db.Collection("student")
    stu := student{"张三",3}
    insertOneResult, err := collection.InsertOne(context.TODO(),stu)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println("insert id: ",insertOneResult.InsertedID) // insert id:  ObjectID("638ac72336193f88b47d19be")
    // 这个 id 在业务里怎么提取使用呢？
}
func insertMany(db *mongo.Database) {
    s1 := student{"Alice",18}
    s2 := student{"Bob",19}
    students := []interface{}{s1,s2}
    collection := db.Collection("student")
    insertManyRes, err := collection.InsertMany(context.TODO(),students) // insert ids:  [ObjectID("638ac7f1266dc9480264ce62") ObjectID("638ac7f1266dc9480264ce63")]
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println("insert ids: ", insertManyRes.InsertedIDs)
}
func update(db *mongo.Database) {
    filter := bson.D{{"name", "张三"}}
    update := bson.D{
        {"$inc", bson.D{
            {"age",200},
        }},
    }
    collection := db.Collection("student")
    updateRes, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println(updateRes.MatchedCount, updateRes.ModifiedCount, updateRes.UpsertedCount)
}

func showSingleRes(db *mongo.Database) {
    collection := db.Collection("student")
    filter := bson.D{{"name", "小王子"}}
    singleResult := collection.FindOne(context.TODO(), filter)
    var res student
    err := singleResult.Decode(&res)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println(res)
}

func showMany(db *mongo.Database) {
    // 设置最多返回 2 个
    findOptions := options.Find() // 这个用法需要再看
    findOptions.SetLimit(2)

    collection := db.Collection("student")
    cursor, err := collection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        log.Println(err)
        return
    }

    var res []*student
    for cursor.Next(context.TODO()) {
        var stu student
        err := cursor.Decode(&stu)
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Printf("stu stu:%+v, %v\n", stu, stu) // {Name:小王子 Age:18}, {小王子 18}
        res = append(res, &stu)
    }

    if err := cursor.Err(); err != nil {
        log.Fatal(err)
        return
    }
    cursor.Close(context.TODO())
    //fmt.Printf("res: %#v\n",res)
}

func delete(db *mongo.Database) {
    collection := db.Collection("student")
    filter := bson.D{{"name", "Bob"}}
    deleteResult1, err := collection.DeleteOne(context.TODO(),filter)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println(deleteResult1.DeletedCount)
    // 删除所有
    deleteResult2, err := collection.DeleteMany(context.TODO(),bson.D{})
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println(deleteResult2.DeletedCount)
}

type student struct {
    Name string  // 字段必须大写，外部包才能在查询时修改传进去的对象字段
    Age int
}

func main() {
    Connect()
    defer Close()
    db := getDB("jinyun")
    //showSingleRes(db)
    //insertOne(db)
    //insertMany(db)
    //delete(db)
    showMany(db)


}
