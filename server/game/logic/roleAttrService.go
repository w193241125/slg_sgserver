package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/model/data"
)

var RoleAttrService = &roleAttrService{}

type roleAttrService struct {
}

func (r *roleAttrService) TryCreate(rid int, conn net.WSConn) error {
	//根据用户ID查询对应的游戏角色..
	role := &data.RoleAttribute{}
	get, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询数据库rid出错")
	}
	if !get {
		role.Id = rid
		role.UnionId = 0
		role.ParentId = 0
		_, err := db.Engine.Table(role).Insert(role)
		if err != nil {
			log.Println("插入初始角色属性失败", err)
			return common.New(constant.DBError, "插入初始角色属性失败")
		}
	}
	return nil
}
