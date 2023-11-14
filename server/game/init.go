package game

import (
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/game/controller"
	"sgserver/server/game/gameConfig"
)

var Router = &net.Router{}

func Init() {
	db.TestDB()
	//加载基础配置
	gameConfig.Base.Load()
	//加载地图配置
	gameConfig.MapBuildConf.Load()
	InitRouter()
}

func InitRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
}
