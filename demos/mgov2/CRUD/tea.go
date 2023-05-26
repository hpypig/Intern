package CRUD

import "time"

func parseOpLogAndSync2Elasitc() *e.Error {

    lastTime, errs := config.TsWorker.ReadLastTime()
    if errs != nil {
        log.Errorf("ReadLastTime err:%s", errs)
        return errs
    }

    cond := bson.M{"ts": bson.M{"$gte": lastTime}}

    session, errs := config.MongoPool.Get()
    if errs != nil {
        log.Errorf("conig.MongoPool.Get err:%s", errs)
        return errs
    }
    defer config.MongoPool.Put(session, false)

    insertCnt := 0
    updateCnt := 0
    deleteCnt := 0
    log.Tracef("begin to tail, lastTime:%d", lastTime)

    iter := session.DB("local").C("oplog.rs").Find(cond).Sort("$natural").Tail(time.Duration(config.TsWorker.TsItem.TailTimeout) * time.Second)
    log.Tracef("get iter ok, lastTime:%d, %d", lastTime, lastTime>>32)

    for {
        oplog := model.NewOpLog()

        for iter.Next(oplog) {

            //上报流水统计
            reportOplog(oplog)

            lastTime = oplog.Ts
            _, present := config.Global.Ns[oplog.Ns]

            if !present {
                continue
            }
            config.Statsd.IncrTotal(model.STATS_MGO2ES_ARTICLE, 1)

            if oplog.Op == "i" {
                insertCnt += 1
            } else if oplog.Op == "u" {
                updateCnt += 1
            } else if oplog.Op == "d" {
                deleteCnt += 1
            }

            workChan <- oplog

            //log.Debugf("ns:%s, op:%s, o2:%+v", oplog.Ns, oplog.Op, oplog.O2.GetId())
            if (insertCnt+updateCnt)%config.TsWorker.TsItem.WindowSize == 0 {
                config.TsWorker.UpdateLastTime(lastTime)
            }
            oplog = model.NewOpLog()
            config.Statsd.SetStatus(model.STATS_MGO2ES_ARTICLE, int(ConvertTs2YmdHms(lastTime)))
        }

        iterErr := iter.Err()
        if iterErr != nil {
            iter.Close()
            log.Errorf("iterErr:%s, lastTime:%d", iterErr, lastTime)
            config.TsWorker.UpdateLastTime(lastTime)
            return e.NewErrorWithError(1, iterErr)
        }

        if iter.Timeout() {
            //log.Debugf("timeout")
            config.TsWorker.UpdateLastTime(lastTime)
            continue
        }

        config.TsWorker.UpdateLastTime(lastTime)

        log.Tracef("new tail:%d", lastTime)
        cond := bson.M{"ts": bson.M{"$gte": lastTime}}
        iter = session.DB("local").C("oplog.rs").Find(cond).Sort("$natural").Tail(time.Duration(config.TsWorker.TsItem.TailTimeout) * time.Second)
    }
}

