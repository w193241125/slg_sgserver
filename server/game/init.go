package game

import (
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/game/controller"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/gameConfig/general"
	"sgserver/server/game/logic"
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
	//加载武将信息
	general.General.Load()
	//记载技能配置
	gameConfig.Skill.Load()

	//初始化前加载信息
	logic.BeforeInit()
	//加载所有建筑信息
	logic.RoleBuildService.Load()
	//加载所有城池信息
	logic.RoleCityService.Load()
	//加载所有角色属性
	logic.RoleAttrService.Load()

	InitRouter()
}

func InitRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
	controller.DefaultGeneralController.Router(Router)
	controller.DefalArmyController.Router(Router)
	controller.DefaultWarController.Router(Router)
	controller.DefaultSkillController.Router(Router)
	controller.InteriorController.Router(Router)
}
