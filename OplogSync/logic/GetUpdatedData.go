package logic

import (
    "gopkg.in/mgo.v2/bson"
    "oplogsync/dao/mongo"
    "oplogsync/entities"
)

// GetUpdatedData 获取 oplog 并分发 data
func GetUpdatedData(oplog entities.Oplog) *entities.UpdatedDataResponse {


        var _id bson.ObjectId
        if oplog.Op == "i" {
            _id = oplog.O.Id
        } else if oplog.Op == "u"{
            _id = oplog.O2.Id
        } else {
            return nil
        }
        //_id := oplog.O.Id.Hex()

        //-----------------
        // 从 mongo 找到数据
        var content entities.NewsContent
        //mongo.FindContentByObjectID(bson.ObjectIdHex(_id), &content) // dao
        mongo.FindContentByObjectID(_id, &content)
        updatedDate := entities.UpdatedDataResponse {
            Op: oplog.Op,
            Id: content.Id,
            Data: content,
        }
        //log.Printf("GetUpdatedData: %+v %+v\n", updatedDate.Op, updatedDate.Id)  // 有的 id 是空的，因为资讯已经被删除了
        //log.Printf("GetUpdatedData: %+v\n", updatedDate)
        return &updatedDate


}
