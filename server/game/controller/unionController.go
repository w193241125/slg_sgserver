package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
)

var UnionController = &uinonController{}

type uinonController struct {
}

func (u uinonController) Router(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", u.list, middleware.CheckRole())
}

func (u *uinonController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.ListRsp{}
	//查询所有联盟
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	uns, err := logic.CoalitionService.List()
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = uns
	rsp.Body.Msg = rspObj
}
