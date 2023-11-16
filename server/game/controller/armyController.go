package controller

import (
	"github.com/mitchellh/mapstructure"
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var DefalArmyController = &ArmyController{}

type ArmyController struct {
}

func (a *ArmyController) Router(router *net.Router) {
	g := router.Group("army")
	g.Use(middleware.Log())
	g.AddRouter("myList", a.mylist, middleware.CheckRole())
}

func (a *ArmyController) mylist(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ArmyListReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ArmyListRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	role, _ := req.Conn.GetProperty("role")
	r := role.(*data.RoleModel)
	arms, err := logic.ArmyService.GetArmysByCity(r.RId, reqObj.CityId)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = arms
	rspObj.CityId = reqObj.CityId
	rsp.Body.Msg = rspObj
}
