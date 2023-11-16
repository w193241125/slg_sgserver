package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefaultGeneralController = &GeneralController{}

type GeneralController struct {
}

func (gc *GeneralController) Router(router *net.Router) {
	g := router.Group("general")
	g.Use(middleware.Log())
	g.AddRouter("myGenerals", gc.myGenerals, middleware.CheckRole())
}

func (gc *GeneralController) myGenerals(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//角色拥有的武将,查询出来即可,
	//初始化进入游戏的时候, 没有武将,就随机2个武将给他.
	rspObj := &model.MyGeneralRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.RoleModel).RId
	gs, err := logic.GeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Generals = gs
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
