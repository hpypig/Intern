





















































































































































db.createCollection('UpdatedContent');
// 更新
db.getCollection('UpdatedContent').insertOne({
    "op": "u",                  // i:新增，u:更新
    "id": "5ab50dc1c70e4127558d4ece",   //资讯ID
    "data":{
        "id" : "63aace6cbc488a00077fe943",
        "title" : "媒体打了港股甲标签，深擎打了SZ乙标签",
        "subtitle" : "",
        "content" : "<p>财联社6月16日讯（记者 徐昊）腾势D9的到来，再次点燃了比亚迪在高端市场的希望。</p><p>6月16日，腾势销售事业部总经理赵长江在社交平台上披露了腾势D9的订单情况和新车交付计划。</p><p>“大家在没有看到车和我们没有花广告去曝光的情况下，腾势自5月16号发布预售以来订单已破2万台。”赵长江的“凡尔赛”引发了网友的热议，有网友称“无比期待腾势首款MPV”。</p>",
        "roundup" : "",
        "categories" : [],
        "columns" : [
            "5cbd6668cc9d755ecdc6fb1a"
        ],
        "authors" : [
            {
                "name" : "巢锡镇"
            }
        ],
        "media" : "财联社",
        "source" : 1004,
        "picUrl" : "https://img.cls.cn/images/20220616/v5TQ2yvlCj.jpg",
        "privilege" : 0,
        "createTime" : NumberLong(1655361117375),
        "updateTime" : NumberLong(1655430841048),
        "publishTime" : NumberLong(1655361117375),
        "offlineTime" : NumberLong(0),
        "status" : 40,
        "stocks" : [
            {
                "name" : "比亚迪",
                "market" : "SZ",
                "code" : "002594"
            },
            {
                "name" : "银河娱乐",
                "market" : "HK",
                "code" : "00027"
            }
        ]
    }
})
//-----------------------------------
db.getCollection('NewsContent').find({"id":"testnews1"});
// 插入测试数据，更新，然后查看oplog
db.getCollection('NewsContent').insertOne({
    "source": 1007,
    "id": "testnews1",
    "title": "测试要闻",
    "subtitle": "",
    "media": "e公司",
    "content": "测试内容----------------更改前",
    "privilege": 0,
    "status": 40,
    "dataFlag": 0,
    "createTime": 1670327774,
    "updateTime": 1670327776,
    "publishTime": 1670327583,
    "categories": [

    ],
    "columns": [
        "5cbd6668cc9d755ecdc6fb20"
    ],
    "picUrl": "",
    "bigPicUrl": "",
    "roundup": "",
    "keywords": [

    ],
    "stocks": [
        {
            "market": "SH",
            "code": "600309",
            "name": "万华化学",
            "tag": {
                "weight": 100,
                "emotion": "中性",
                "emotionWeight": 100
            }
        },
        {
            "market": "SH",
            "code": "600077",
            "name": "宋都股份",
            "tag": {
                "weight": 78.38999938964844,
                "emotion": "正面",
                "emotionWeight": 100
            }
        }
    ],
    "sourceName": "证券时报",
    "extPrivilege": "",
    "columnsObj": [
        {
            "id": "5cbd6668cc9d755ecdc6fb20",
            "title": "证券时报"
        }
    ],
    "txtType": 1,
    "offlineTime": 0,
    "uuid": "1007_secu_time_748428",
    "disclaimer": "投资有风险，外部资讯仅供参考，不代表广发证券股份有限公司对其内容的认可或推荐，不构成广发证券股份有限公司做出的投资建议或对任何证券投资价值观点的认可。投资者应当自主进行投资决策，对投资者因依赖上述信息进行投资决策而导致的财产损失，本公司不承担法律责任。",
    "external": 2,
    "links": [
        {
            "word": "国华",
            "type": 1,
            "target": "HK_00370"
        },
        {
            "word": "东方电气",
            "type": 1,
            "target": "HK_01072"
        },
        {
            "word": "万华化学",
            "type": 1,
            "target": "SH_600309"
        },
        {
            "word": "隆基绿能",
            "type": 1,
            "target": "SH_601012"
        },
        {
            "word": "明阳智能",
            "type": 1,
            "target": "SH_601615"
        },
        {
            "word": "光大证券",
            "type": 1,
            "target": "SH_601788"
        },
        {
            "word": "天合光能",
            "type": 1,
            "target": "SH_688599"
        },
        {
            "word": "资源优势",
            "type": 1,
            "target": "SZ_399319"
        },
        {
            "word": "环渤海",
            "type": 1,
            "target": "SZ_399357"
        },
        {
            "word": "能源行业",
            "type": 3,
            "target": "510610"
        }
    ],
    "tts": {
        "url": "https://tts-1251438792.file.myqcloud.com/v100/202212/638f2ddee4f89200073157fa-u4f4er2n.mp3?sign=771bf232d4de76d30ba8f5037528438f&t=1670489833",
        "duration": 559010,
        "size": 1677285
    },
    "maskTitle": "万华化学加速新能源领域布局 携手国能集团等建海上光伏发电项目",
    "listenCount": 0,
    "isOriginal": false,
    "readCnt": 336,
    "relativeProductType": 0,
    "contentLength": 0,
    "likeCnt": 0,
    "dislikeCnt": 0,
    "madId": "1fRt",
    "plazaId": "1050zq"
});




db.getCollection('NewsContent').update(
    {"id":"testnews1"},
    {
       $set:{"content": "测试内容----更改后"}
    }
);

db.NewsContent.findOne({"id":"testnews1"})
use local
// 这样其实无法查到最新数据（更新逻辑需要思考）
db.oplog.rs.findOne({"o.id":"testnews1"})
//db.getCollection('oplog.rs').find({"id":"testnews1"})
