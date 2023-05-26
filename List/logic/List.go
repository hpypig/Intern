package logic

import (
    "fmt"
    "github.com/hpypig/Intern/dao/mongo"
    "github.com/hpypig/Intern/dao/redis"
    "github.com/hpypig/Intern/entities"
    "strings"
    "time"
)

func GetInfo(p *entities.ParamStockNews) (data []entities.BriefInfo, err error) {
    //fmt.Printf("%+v\n",*p)
    typeMarketCode := fmt.Sprintf("%d_%s",p.TxtType,p.MarketCode)
    idsKey := redis.StockListKeyPrefix + "_" + strings.ToUpper(typeMarketCode)
    start := (p.Page - 1) * p.PageSize
    end := p.Page * p.PageSize - 1

    //fmt.Printf("idsKey:%s, start:%d, end:%d\n",idsKey, start, end)
    //t1 := time.Now()
    // 获取 id 列表
    idTitleMedia, err := redis.GetIdsByRevRange(idsKey, start, end) // 按高分（时间大的）获取列表
    //t2 := time.Now()
    if err != nil {
        //log.Println("List-GetInfo err:", err)
        fmt.Println("List-GetInfo err:", err)
        return nil, err
    }
    //fmt.Println("List-GetInfo redis-mongo")
    for _,v:= range idTitleMedia {
        strs := strings.Split(v,"_")
        if len(strs) > 3 {
            fmt.Println("err!!! idTitleMedia length > 3: ", v)
        }
        data = append(data,entities.BriefInfo{Id: strs[0], Title: strs[1], Media: strs[2]})
    }
    //t3 := time.Now()
    //fmt.Println("Get Info 时间消耗：", t2.Sub(t1), t3.Sub(t2))
    return data, nil
}

func GetImportantNewsInfo(limit int) (data []entities.BriefInfo) {
    // 从 redis 获取 30 条索引
    key := redis.ImportantNewsListKey
    var start, end int64
    end = int64(limit)-1
    idTitleMedia, _ := redis.GetIdsByRevRange(key, start, end)
    // 获取标题、时间、媒体
    for _,v:= range idTitleMedia {
        strs := strings.Split(v,"_")
        if len(strs) > 3 {
            fmt.Println("err!!! idTitleMedia length > 3: ", v)
        }
        data = append(data,entities.BriefInfo{Id: strs[0], Title: strs[1], Media: strs[2]})
    }
    return
    // 调用 listgen 请求实时推送，不是在这儿做的，是在 websocket 接口做的
}

func GetNewsContent(id string, data *entities.NewsContent)  {
    mongo.GetNewsContentById(id, data)
}

func GetColumnInfo(p *entities.ParamColumnRequest) (data []entities.BriefInfo){
    key := redis.ColumnListKeyPrefix + "_" + p.Id
    start := (p.Page - 1) * p.PageSize
    end := p.Page * p.PageSize - 1

    idTitleMedia, _ := redis.GetIdsByRevRange(key, start, end)
    // 从 mongodb 获取标题、时间、媒体
    // 获取标题、时间、媒体
    for _,v:= range idTitleMedia {
        strs := strings.Split(v,"_")
        if len(strs) > 3 {
            fmt.Println("err!!! idTitleMedia length > 3: ", v)
        }
        data = append(data,entities.BriefInfo{Id: strs[0], Title: strs[1], Media: strs[2]})
    }
    return
}


//2.3 前未改redis value格式的版本-------------

func GetInfo2(p *entities.ParamStockNews) (data []entities.BriefInfo, err error) {
    //fmt.Printf("%+v\n",*p)
    typeMarketCode := fmt.Sprintf("%d_%s",p.TxtType,p.MarketCode)
    idsKey := redis.StockListKeyPrefix + "_" + strings.ToUpper(typeMarketCode)
    start := (p.Page - 1) * p.PageSize
    end := p.Page * p.PageSize - 1

    //fmt.Printf("idsKey:%s, start:%d, end:%d\n",idsKey, start, end)
    //t1 := time.Now()
    // 获取 id 列表
    ids, err := redis.GetIdsByRevRange(idsKey, start, end) // 按高分（时间大的）获取列表
    //t2 := time.Now()
    if err != nil {
        //log.Println("List-GetInfo err:", err)
        fmt.Println("List-GetInfo err:", err)
        return nil, err
    }
    //fmt.Println("List-GetInfo redis-mongo")

    t3 := time.Now()
    data = mongo.GetBriefInfoByIds2(ids)
    //fmt.Println("GetInfo 时间消耗:",t2.Sub(t1), " ",time.Since(t3))
    fmt.Println("GetInfo mongo 时间消耗:", time.Since(t3))
    return data, nil
    // 这里有些问题，上面循环里的 error 没处理
}



func GetImportantNewsInfo2(limit int) (data []entities.BriefInfo) {
    // 从 redis 获取 30 条索引
    key := redis.ImportantNewsListKey
    var start, end int64
    end = int64(limit)-1
    ids, _ := redis.GetIdsByRevRange(key, start, end)
    // 从 mongodb 获取标题、时间、媒体
    data = mongo.GetBriefInfoByIds(ids)
    return
    // 调用 listgen 请求实时推送，不是在这儿做的，是在 websocket 接口做的
}

func GetNewsContent2(id string, data *entities.NewsContent)  {
    mongo.GetNewsContentById(id, data)
}

func GetColumnInfo2(p *entities.ParamColumnRequest) (data []entities.BriefInfo){
    key := redis.ColumnListKeyPrefix + "_" + p.Id
    start := (p.Page - 1) * p.PageSize
    end := p.Page * p.PageSize - 1

    ids, _ := redis.GetIdsByRevRange(key, start, end)
    // 从 mongodb 获取标题、时间、媒体
    data = mongo.GetBriefInfoByIds(ids)
    return
}
