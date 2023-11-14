package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/game/gameConfig"
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
	//根据用户ID查询对应的游戏角色..
	role := &data.Role{}
	get, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		rsp.Body.Code = constant.DBError
		return
	}
	if get {
		rsp.Body.Code = constant.OK
		rsp.Body.Msg = rspObj
		rid := role.RId
		roleRes := &data.RoleRes{}
		ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源出错", err)
			rsp.Body.Code = constant.DBError
			return
		}
		if !ok {
			roleRes.RId = rid
			roleRes.Gold = gameConfig.Base.Role.Gold
			roleRes.Decree = gameConfig.Base.Role.Decree
			roleRes.Grain = gameConfig.Base.Role.Grain
			roleRes.Iron = gameConfig.Base.Role.Iron
			roleRes.Stone = gameConfig.Base.Role.Stone
			roleRes.Wood = gameConfig.Base.Role.Wood
			_, err := db.Engine.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源错误", err)
				rsp.Body.Code = constant.DBError
				return
			}
		}
		rspObj.RoleRes = roleRes.ToModel().(model.RoleRes)
		rspObj.Role = role.ToModel().(model.Role)
		token, err := utils.Award(16)
		if err != nil {
			return
		}
		rspObj.Time = time.Now().UnixNano() / 1e6
		rspObj.Token = token

	} else {
		log.Println("无角色", err)
		rsp.Body.Code = constant.DBError
		return
	}
	//根据角色ID 查询角色拥有的资源, 有就返回, 没有就初始化.
}
