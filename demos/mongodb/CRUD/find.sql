/*
ref:
    https://www.liwenzhou.com/posts/Go/mongodb/#:~:text=()%0Atrue-,%E6%96%87%E6%A1%A3%E5%B8%B8%E7%94%A8%E5%91%BD%E4%BB%A4,-%E6%8F%92%E5%85%A5%E4%B8%80%E6%9D%A1%E6%96%87
    https://www.jianshu.com/p/dbf965f8d314
*/
use gf
// 查询所有
db.getCollection("StockNodes").find();
db.StockNodes.find(
    {id:"638d9e50e4f8920007314bac"}
)

db.NewsContent.find(
    {id:"638d9e50e4f8920007314bac"}
)
// $in 用法
db.NewsContent.find(
    {id:{$in:["638d9e50e4f8920007314bac","638f2ddee4f89200073157fa","638f0c2be4f8920007315685"]}}
)

db.oplog.rs.find({ts:{$gt:1670733911},ns:"gf.NewsContent"})
.sort({field:1}) //升序   -1 降序
