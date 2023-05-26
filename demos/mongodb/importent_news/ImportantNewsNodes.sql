use gf

// db.getCollection("ImportantNewsNode").drop()
// 要闻节点（索引）
db.ImportantNewsNodes.find({id:"638f2ddee4f89200073157fa"})
db.createCollection('ImportantNewsNodes');
db.getCollection('ImportantNewsNodes').insertOne({
    "id": "638f2ddee4f89200073157fa",
    "title": "万华化学加速新能源领域布局 携手国能集团等建海上光伏发电项目",
    "picUrl": "",
    "media": "e公司",
    "stocks": [
        {
            "market": "SH",
            "code": "600309",
            "name": "万华化学"
        },
        {
            "market": "SH",
            "code": "600077",
            "name": "宋都股份"
        }
    ],
    "privilege": 0,
    "publishTime": 1670327583000,
    "isTop": false,
    "finHeadLineType": 0,
    "linkUrl": "",
    "outerId": "",
    "readCnt": 112,
    "manualRank": 10000,
    "manualTags": ""
})
