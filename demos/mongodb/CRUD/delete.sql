// 和 find 格式应该差不多，{过滤},{指定显示}
db.getCollection("xxx").deleteOne({})
db.getCollection("xxx").deleteMany({},{})
// 删除所有
db.getCollection("xxx").deleteMany({})
