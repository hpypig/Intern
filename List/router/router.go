package router

import (
    "github.com/gin-contrib/pprof"
    "github.com/gorilla/websocket"
    "github.com/hpypig/Intern/midware"
    "github.com/hpypig/Intern/ws"

    //"fmt"
    "github.com/gin-gonic/gin"
    "github.com/hpypig/Intern/entities"
    "github.com/hpypig/Intern/logic"
    "log"
    "net/http"
)

var Hub *ws.Hub  // websocket 客户端的注册和广播

func Init(r *gin.Engine) {
    pprof.Register(r)

    Hub = ws.NewHub()
    go Hub.Run()

    r.Use(midware.Cors())

    r.LoadHTMLGlob("./static/*")

    r.GET("/", ServeHome)
    // 根据 market_code page page_size txtType 获取资讯标题、时间、媒体等
    r.GET("/list/article/stocks/batch", GetStockInfoHandler)
    // 获取最新的 30 条要闻标题
    r.GET("/ImportantNews",GetImportantNewsInfoHandler)
    // 获取栏目标题
    r.GET("/list/article/column",GetColumnInfoHandler)           // 最后打斜线和不打有区别吗？加了斜线报301错误，前端报跨域错误，可是明明加了cors
    // 获取具体内容
    r.GET("read/article/:id",GetNewsContentHandler)

}

func ServeHome(c *gin.Context) {
    c.HTML(http.StatusOK,"index.html",nil)
}

func GetStockInfoHandler(c *gin.Context) {
    //fmt.Println("in GetSimpleInfoHandler")
    p := entities.NewParamStockNews()
    if err := c.ShouldBindQuery(p); err != nil {
        log.Printf("GetStockInfoHandler:get params err: %v \n", err)
        return
    }
    //fmt.Printf("GetSimpleInfoHandler:params: %+v\n", p)

    //num, err := strconv.ParseInt(txtType,10,64)
    //if err != nil {
    //    log.Println("GetSimpleInfo parseInt err: ",err)
    //    return
    //}

    infoList, err := logic.GetInfo(p)
    //fmt.Printf("infoList: %+v\n",infoList)
    if err != nil {
        log.Println("GetSimpleInfoHandler err: ", err)
        return
    }
    //fmt.Println("infoList len: ", len(infoList)) // 只有20条，应该是 iter 的时候少遍历了 1 条
    // 转换为字节
    c.JSON(http.StatusOK, &entities.ResponseData{
        ErrCode: 0,
        ErrMsg: "",
        Data:infoList,
    })
}



func GetImportantNewsInfoHandler(c *gin.Context) {
    var upgrader websocket.Upgrader
    upgrader.CheckOrigin = func(r *http.Request) bool {
        return true
    }

    limit := 30
    infoList := logic.GetImportantNewsInfo(limit)

    if !c.IsWebsocket() {
        c.JSON(http.StatusOK, &entities.ResponseData{
            ErrCode: 0,
            ErrMsg: "",
            Data:infoList,
        })
    } else {
        // 升级、注册、运行
        //conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println("upgrade err: ",err)
            return
        }
        client := ws.NewClient(Hub, conn, make(chan []byte, 256))
        ws.PushImportantNews(client, infoList)

    }


}

func GetColumnInfoHandler(c *gin.Context) {
    p := entities.NewParamColumnRequest()
    err := c.ShouldBindQuery(p)
    if err != nil {
        log.Println("GetColumnInfoHandler err:", err)
        return
    }
    infoList := logic.GetColumnInfo(p)
    c.JSON(http.StatusOK, &entities.ResponseData {
        ErrCode: 0,
        ErrMsg: "",
        Data:infoList,
    })
}

func GetNewsContentHandler(c *gin.Context) {
    id := c.Param("id")
    var newsContent entities.NewsContent
    logic.GetNewsContent(id, &newsContent)
    c.JSON(http.StatusOK, &entities.ResponseData{
        ErrCode: 0,
        ErrMsg: "",
        Data: newsContent,
    })

}
