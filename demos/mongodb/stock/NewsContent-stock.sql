use gf;
db.createCollection('NewsContent');


// 个股 之后还要把另外两个类别加进来
db.getCollection('NewsContent').insertOne({
    "source": 10004,
    "id": "638d9e50e4f8920007314bac",
    "title": "券商股震荡走高 国联证券冲击涨停",
    "subtitle": "",
    "media": "财联社直播",
    "content": "国联证券冲击涨停，广发证券、东方财富、中信建投、国金证券、中信证券等跟涨。",
    "privilege": 0,
    "status": 40,
    "dataFlag": 0,
    "createTime": 1670225488,
    "updateTime": 1670225488,
    "publishTime": 1670205452,
    "categories": [

    ],
    "columns": [
        "5cbd6668cc9d755ecdc6fb25"
    ],
    "picUrl": "",
    "bigPicUrl": "",
    "roundup": "",
    "keywords": [

    ],
    "stocks": [
        {
            "market": "SH",
            "code": "601456",
            "name": "国联证券",
            "tag": {
                "weight": 98.26000213623047,
                "emotion": "正面",
                "emotionWeight": 100
            }
        },
        {
            "market": "SZ",
            "code": "000776",
            "name": "广发证券",
            "tag": {
                "weight": 87.0999984741211,
                "emotion": "正面",
                "emotionWeight": 100
            }
        },
        {
            "market": "SZ",
            "code": "300059",
            "name": "东方财富",
            "tag": {
                "weight": 83.05000305175781,
                "emotion": "正面",
                "emotionWeight": 100
            }
        },
        {
            "market": "SH",
            "code": "601066",
            "name": "中信建投",
            "tag": {
                "weight": 80.94000244140625,
                "emotion": "正面",
                "emotionWeight": 100
            }
        },
        {
            "market": "SH",
            "code": "600109",
            "name": "国金证券",
            "tag": {
                "weight": 78.77999877929688,
                "emotion": "正面",
                "emotionWeight": 100
            }
        },
        {
            "market": "SH",
            "code": "600030",
            "name": "中信证券",
            "tag": {
                "weight": 76.55999755859375,
                "emotion": "正面",
                "emotionWeight": 100
            }
        }
    ],
    "sourceName": "财联社",
    "extPrivilege": "",
    "columnsObj": [
        {
            "id": "5cbd6668cc9d755ecdc6fb25",
            "title": "财联社直播"
        }
    ],
    "txtType": 1,
    "offlineTime": 0,
    "uuid": "10004_cailianshe_live_1202807",
    "disclaimer": "投资有风险，外部资讯仅供参考，不代表广发证券股份有限公司对其内容的认可或推荐，不构成广发证券股份有限公司做出的投资建议或对任何证券投资价值观点的认可。投资者应当自主进行投资决策，对投资者因依赖上述信息进行投资决策而导致的财产损失，本公司不承担法律责任。",
    "external": 2,
    "links": [
        {
            "word": "国联证券",
            "type": 1,
            "target": "HK_01456"
        },
        {
            "word": "中信证券",
            "type": 1,
            "target": "HK_06030"
        },
        {
            "word": "国金证券",
            "type": 1,
            "target": "SH_600109"
        },
        {
            "word": "中信建投",
            "type": 1,
            "target": "SH_601066"
        },
        {
            "word": "广发证券",
            "type": 1,
            "target": "SZ_000776"
        },
        {
            "word": "东方财富",
            "type": 1,
            "target": "SZ_300059"
        }
    ],
    "tts": {
        "url": "https://tts-1251438792.file.myqcloud.com/v100/202212/638d9e50e4f8920007314bac-9bmt26wm.mp3?sign=491079be9d18f15cd9281976b77ab6d2&t=1670489680",
        "duration": 12670,
        "size": 38277
    },
    "maskTitle": "券商股震荡走高 国联证券冲击涨停",
    "listenCount": 0,
    "isOriginal": false,
    "readCnt": 1286,
    "relativeProductType": 0,
    "contentLength": 0,
    "industries": [
        {
            "market": "SW21",
            "code": "490000",
            "name": "非银金融",
            "tag": {
                "weight": 100,
                "emotion": "",
                "emotionWeight": 0
            }
        }
    ],
    "likeCnt": 0,
    "dislikeCnt": 0,
    "madId": "1dkq"
})
