package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefaultWarController = &WarController{}

type WarController struct {
}

func (w *WarController) Router(router *net.Router) {
	g := router.Group("war")
	g.AddRouter("report", w.report)
}

func (w *WarController) report(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查找战报表 得出数据
	rspObj := &model.WarReportRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.RoleModel).RId
	reports, err := logic.WarService.GetWarReport(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = reports
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
}
