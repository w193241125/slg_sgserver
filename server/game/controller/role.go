package controller

import (
	"github.com/mitchellh/mapstructure"
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/model"
	"sgserver/utils"
)

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//进入游戏逻辑
	// Session 是否合法, 合法可取出登录用户的ID

	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		rsp.Body.Code = constant.InvalidParam
		return
	}
	session := reqObj.Session
	_, claim, err := utils.ParseToken(session)
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	uid := claim.Uid
	err = logic.RoleService.EnterServer(uid, rspObj, req.Conn)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
	//根据角色ID 查询角色拥有的资源, 有就返回, 没有就初始化.
}
