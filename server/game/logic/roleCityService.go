package logic

import (
	"log"
	"math/rand"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/global"
	"sgserver/server/game/model/data"
	"time"
)

var RoleCityService = &roleCityService{}

type roleCityService struct {
}

func (r *roleCityService) InitCity(rid int, nickname string, conn net.WSConn) error {
	//根据用户ID查询对应的游戏角色..
	roleCity := &data.RoleCity{}
	get, err := db.Engine.Table(roleCity).Where("rid=?", rid).Get(roleCity)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return common.New(constant.DBError, "查询数据库rid出错")
	}
	if !get {
		//初始化
		roleCity.X = rand.Intn(global.MapWith)
		roleCity.X = rand.Intn(global.MapHeight)
		//城池是否能在这个坐标创建, 需要判断. 五个之内不能有玩家.
		roleCity.RId = rid
		roleCity.Name = nickname
		roleCity.CurDurable = gameConfig.Base.City.Durable
		roleCity.CreatedAt = time.Now()
		roleCity.IsMain = 1
		_, err := db.Engine.Table(roleCity).Insert(roleCity)
		if err != nil {
			log.Println("城池初始化失败", err)
			return common.New(constant.DBError, "插入城池初始化信息失败")
		}
	}
	return nil
}
