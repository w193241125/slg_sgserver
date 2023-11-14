package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"time"
)

var RoleService = &roleService{}

type roleService struct {
}

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, conn net.WSConn) error {
	//根据用户ID查询对应的游戏角色..
	role := &data.Role{}
	get, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询数据库uid出错")
	}
	if get {
		rid := role.RId
		roleRes := &data.RoleRes{}
		ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源出错", err)
			return common.New(constant.DBError, "查询数据库rid出错")
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
				return common.New(constant.DBError, "插入角色资源错误")
			}
		}
		rsp.RoleRes = roleRes.ToModel().(model.RoleRes)
		rsp.Role = role.ToModel().(model.Role)
		token, _ := utils.Award(16)
		rsp.Time = time.Now().UnixNano() / 1e6
		rsp.Token = token
		//将角色信息存入socket中
		conn.SetProperty("role", role)
		//初始化玩家属性
		if err := RoleAttrService.TryCreate(rid, conn); err != nil {
			return common.New(constant.DBError, "尝试创角失败")
		}
	} else {
		log.Println("无角色", err)
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	return nil
}
