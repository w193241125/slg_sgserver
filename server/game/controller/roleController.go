package controller

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"time"
)

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.Use(middleware.Log())
	g.AddRouter("create", r.create)
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList)
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
	err = logic.RoleService.EnterServer(uid, rspObj, req)
	if err != nil {
		rspObj.Time = time.Now().UnixNano() / 1e6
		rsp.Body.Msg = rspObj
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
	//根据角色ID 查询角色拥有的资源, 有就返回, 没有就初始化.
}

func (r *RoleController) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//分别根据角色ID查询将军  资源 建筑 城池 武将
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name

	rid := role.(*data.RoleModel).RId
	rspObj := &model.MyRolePropertyRsp{}
	//查询资源
	rspObj.RoleRes, err = logic.RoleService.GetRoleRes(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//城池
	rspObj.Citys, err = logic.RoleCityService.GetRoleCity(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//建筑
	rspObj.MRBuilds, err = logic.RoleBuildService.GetBuilds(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//军队
	rspObj.Armys, err = logic.ArmyService.GetArmys(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//武将
	rspObj.Generals, err = logic.GeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}

	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *RoleController) posTagList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.PosTagListRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	//查询角色属性表
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.RoleModel).RId
	pts, err := logic.RoleAttrService.GetTagList(rid)
	if err != nil {
		rsp.Body.Code = err.(common.MyError).Code()
		return
	}
	rspObj.PosTags = pts
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *RoleController) create(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.CreateRoleRsp{}
	reqObj := &model.CreateRoleReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	cr, _ := json.Marshal(reqObj)
	log.Println(string(cr))
	role := &data.RoleModel{}
	ok, err := db.Engine.Where("uid=?", reqObj.UId).Get(role)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	if ok {
		rsp.Body.Code = constant.RoleAlreadyCreate
		return
	}
	role.UId = reqObj.UId
	role.Sex = reqObj.Sex
	role.NickName = reqObj.NickName
	role.Balance = 0
	role.HeadId = reqObj.HeadId
	role.CreatedAt = time.Now()
	role.LoginTime = time.Now()
	if _, err := db.Engine.Insert(role); err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Role = role.ToModel().(model.Role)
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = role
}
