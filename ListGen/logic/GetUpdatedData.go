package logic

import (
    "ListGen/dao/mongo"
    "ListGen/entities"
    "gopkg.in/mgo.v2/bson"
    "log"
)

// GetUpdatedData 获取 oplog 并分发 data
//func GetUpdatedData(oplogChan <-chan entities.Oplog, updatedDataChans []chan entities.UpdatedDataResponse, flagCh chan struct{}) {
func GetUpdatedData(oplogChan <-chan entities.Oplog, updatedDataChans []chan entities.UpdatedDataResponse, subFlag *SubFlag) {
    for {
        oplog := <-oplogChan
        var _id bson.ObjectId
        if oplog.Op == "i" {
            _id = oplog.O.Id
        } else if oplog.Op == "u"{
            _id = oplog.O2.Id
        } else {
            continue
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
        log.Printf("GetUpdatedData: %+v\n", updatedDate)
        if subFlag.OnPush() {
               log.Printf("GetUpdatedData-1\n")
               updatedDataChans[0] <- updatedDate
               updatedDataChans[1] <- updatedDate
               log.Printf("GetUpdatedData-2\n")
        } else {
               log.Printf("GetUpdatedData-3\n")
               updatedDataChans[0] <- updatedDate // 之后可以考虑异步是否有问题
               log.Printf("GetUpdatedData-4\n")
        }

        // 订阅时 ch 关闭
        //select {
        //case <-flagCh:
        //    log.Printf("GetUpdatedData-1\n")
        //    updatedDataChans[0] <- updatedDate
        //    updatedDataChans[1] <- updatedDate
        //    log.Printf("GetUpdatedData-2\n")
        //default:
        //    log.Printf("GetUpdatedData-3\n")
        //    updatedDataChans[0] <- updatedDate // 之后可以考虑异步是否有问题
        //    log.Printf("GetUpdatedData-4\n")
        //}
    }
}
