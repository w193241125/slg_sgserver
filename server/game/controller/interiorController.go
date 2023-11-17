package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"time"
)

var InteriorController = &interiorController{}

type interiorController struct {
}

func (i *interiorController) Router(router *net.Router) {
	g := router.Group("interior")
	g.Use(middleware.Log())
	g.AddRouter("openCollect", i.openCollect, middleware.CheckRole())
	g.AddRouter("collect", i.collect, middleware.CheckRole())
}

func (i *interiorController) openCollect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.OpenCollectionRsp{}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
	r, _ := req.Conn.GetProperty("role")
	rid := r.(*data.RoleModel).RId
	ra := logic.RoleAttrService.Get(rid)
	if ra != nil {
		rspObj.CurTimes = ra.CollectTimes
		rspObj.Limit = gameConfig.Base.Role.CollectTimesLimit
		//征收间隔时间
		interval := gameConfig.Base.Role.CollectInterval

		//最后征收时间为 0
		if ra.LastCollectTime.IsZero() {
			rspObj.NextTime = 0
		} else {
			if rspObj.CurTimes >= rspObj.Limit {
				//今日征收上限, 下次征收时间是明天(最后一次征收时间为准)
				// 从零点开始可以征收
				y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
				//东八区 time.FixedZone("CST", 8*3600)
				ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
				rspObj.NextTime = ti.UnixNano() / 1e6
			} else {
				ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
				rspObj.NextTime = ti.UnixNano() / 1e6
			}
		}
	}
}

func (i *interiorController) collect(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查询角色资源 得到当前金币
	//查询角色属性 得到征收的相关信息
	// 查询当前 产量, 征收的量是多少
	rspObj := &model.CollectionRsp{}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
	r, _ := req.Conn.GetProperty("role")
	rid := r.(*data.RoleModel).RId
	ra := logic.RoleAttrService.Get(rid)
	if ra == nil {
		rsp.Body.Code = constant.DBError
		return
	}

	//角色资源
	rs := logic.RoleResService.GetRoleRes(rid)
	if rs == nil {
		rsp.Body.Code = constant.DBError
		return
	}
	//产量
	yield := logic.RoleResService.GetYield(rid)
	rs.Gold += yield.Gold
	//进行数据库更新  go channel 一旦需要更新就发更新信号, 消耗方接收消息,进行更新
	rs.SyncExecute()

	rspObj.Gold = yield.Gold
	curTime := time.Now()
	limit := gameConfig.Base.Role.CollectTimesLimit
	interval := gameConfig.Base.Role.CollectInterval
	lastTime := ra.LastCollectTime
	if curTime.YearDay() != lastTime.YearDay() || curTime.Year() != lastTime.Year() {
		ra.CollectTimes = 0
		ra.LastCollectTime = time.Time{}
	}
	//计算征收
	ra.CollectTimes += 1
	ra.LastCollectTime = curTime
	ra.SyncExecute()
	rspObj.Limit = limit
	rspObj.CurTimes = ra.CollectTimes

	if rspObj.CurTimes >= rspObj.Limit {
		//今天已经完成征收了，下一次征收就是第二天(最后一次征收时间为准)
		//第二天 从0点就开始了
		y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
		//东八区 time.FixedZone("CST", 8*3600)
		ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
		rspObj.NextTime = ti.UnixNano() / 1e6
	} else {
		ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
		rspObj.NextTime = ti.UnixNano() / 1e6
	}
}
