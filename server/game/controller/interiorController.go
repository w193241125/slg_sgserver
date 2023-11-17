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
