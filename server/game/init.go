package game

import (
	"sgserver/net"
	"sgserver/server/game/controller"
	"sgserver/server/game/gameConfig"
)

var Router = &net.Router{}

func Init() {
	//加载基础配置
	gameConfig.Base.Load()
	InitRouter()
}

func InitRouter() {
	controller.DefaultRoleController.Router(Router)
}
