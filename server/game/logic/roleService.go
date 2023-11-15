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

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, req *net.WsMsgReq) error {
	//根据用户ID查询对应的游戏角色..
	role := &data.RoleModel{}
	session := db.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		log.Println("事务开启出错", err)
		return common.New(constant.DBError, "数据库出错")
	}

	req.Context.Set("dbSession", session)

	get, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询数据库uid出错")
	}
	if get {
		//已经创角了就去查询角色资源
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
			_, err := session.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源错误", err)
				return common.New(constant.DBError, "插入角色资源错误")
			}
		}
		rsp.RoleRes = roleRes.ToModel().(model.RoleRes)
		rsp.Role = role.ToModel().(model.Role)
		rsp.Time = time.Now().UnixNano() / 1e6
		token, _ := utils.Award(rid)
		rsp.Token = token
		//将角色信息存入socket中
		req.Conn.SetProperty("role", role)
		//初始化玩家属性
		if err := RoleAttrService.TryCreate(rid, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "尝试创角失败")
		}
		//初始化城池

		if err := RoleCityService.InitCity(rid, role.NickName, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "城池初始化失败")
		}
	} else {
		log.Println("无角色,去创角", err)
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	err = session.Commit()
	if err != nil {
		log.Println("事务提交出错")
		return common.New(constant.DBError, "事务提交出错")
	}
	return nil
}

func (r *roleService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	get, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("获取角色资源出错", err)
		return model.RoleRes{}, common.New(constant.DBError, "获取角色资源出错")
	}
	if get {
		return roleRes.ToModel().(model.RoleRes), nil
	}

	return model.RoleRes{}, common.New(constant.RoleNotExist, "角色资源不存在")
}
