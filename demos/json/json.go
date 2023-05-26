package jsont

import (
    "encoding/json"
    "fmt"
)

type Student struct{
    Id int64 `json:"id,string"`   // 注意这里的 ,string
    Name string `json:"name"`
    msg string // 小写的不会被json解析
}

func MaU() {
    // 定义一个结构体切片初始化
    student := Student{123123,"whw","666"}
    // 将结构体转成json格式
    data, err := json.Marshal(student)
    if err == nil{
        // 注意这里将 Id转换为了string
        fmt.Printf("%s\n",data)//{"id":"123123","name":"whw"}
    }

    // json反序列化为结构体 这里的id是 字符串类型的。。。
    s := `{"name":"whw","id":"123123"}`
    var StuObj Student
    if err := json.Unmarshal([]byte(s),&StuObj);err != nil{
        fmt.Println("err>>",err)
    }else{
        // 反序列化后 成了 int64 （,string 的作用）
        fmt.Printf("%T \n",StuObj.Id)// int64
        fmt.Printf("%v \n",StuObj)// {123123 whw }
    }
}
//----------------------------
type Point struct{
    X, Y int
}

type Student2 struct{
    Point
    Name string `json:"__name"`
    Id int64 `json:"id"`
    Age int `json:"-"` // 不解析该字段
    msg string // 小写的不会被json解析
}

func MaU2() {
    // 定义一个结构体切片初始化
    po := Point{1,2}
    student := Student2{po,"whw",123123,12,"阿斯顿发送到发"}
    // 将结构体转成json格式
    data, err := json.Marshal(student)
    if err == nil{
        // 注意这里将 Id转换为了string
        fmt.Printf("%s\n",data)// {"X":1,"Y":2,"__name":"whw","id":123123}
    }
}
type Content struct {
    Source int `json:"source""`
    Content string `json:"content"`
}
type Message struct {
    ErrCode int `json:"errCode"`
    ErrMsg string `json:"errMsg"`
    Data Content `json:"data"`
    //a json.RawMessage
}

// un时不理会mar里面的一些字段可以吗？

func MaU3() {
    body := []byte(`{"errCode":1,"errMsg":"msg","data":{"source":1007,"content":"12344121"},"data2":"123"}`)

    var m Message
    json.Unmarshal(body,&m)
    fmt.Printf("%s\n",body)
    fmt.Printf("%+v\n",m) // {ErrCode:1 ErrMsg:msg Data:{Source:1007 Content:12344121}}
}
