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
	//加载地图单元格配置
	gameConfig.MapRes.Load()
	//加载城池设施配置
	gameConfig.FacilityConf.Load()
	InitRouter()
}

func InitRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
}
