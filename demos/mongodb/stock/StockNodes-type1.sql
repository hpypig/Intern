// db.dropDatabase();
use gf;
db.createCollection('StockNodes');
// show collections;
db.getCollection('StockNodes').find()
db.getCollection('StockNodes').insertOne({
    "id": "638d9e50e4f8920007314bac",
    "title": "券商股震荡走高 国联证券冲击涨停",
    "picUrl": "",
    "publishTime": 1670205452,
    "privilege": 0,
    "media": "财联社直播",
    "stocks": [
        {
            "market": "SZ",
            "code": "000776",
            "name": "广发证券",
            "type": 10
        }
    ],
    "columns": [
        "5cbd6668cc9d755ecdc6fb25"
    ],
    "type": 0,
    "txtType": 1,
    "industries": [
        {
            "market": "SW21",
            "code": "490000",
            "name": "非银金融"
        }
    ],
    "maskTitle": ""
});
