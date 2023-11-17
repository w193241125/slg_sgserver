package logic

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model/data"
)

var RoleResService = &roleResService{}

type roleResService struct {
}

func (r *roleResService) GetYield(rid int) data.Yield {
	// 基础产量+城池设施的产量+建筑产量
	rbYield := RoleBuildService.GetYield(rid)
	cfYield := CityFacilityService.GetYield(rid)
	var y data.Yield
	y.Gold = rbYield.Gold + cfYield.Gold + gameConfig.Base.Role.GoldYield
	y.Stone = rbYield.Stone + cfYield.Stone + gameConfig.Base.Role.StoneYield
	y.Iron = rbYield.Iron + cfYield.Iron + gameConfig.Base.Role.IronYield
	y.Grain = rbYield.Grain + cfYield.Grain + gameConfig.Base.Role.GrainYield
	y.Wood = rbYield.Wood + cfYield.Wood + gameConfig.Base.Role.WoodYield
	return y
}

func (r *roleResService) GetRoleRes(rid int) *data.RoleRes {
	roleRes := &data.RoleRes{}
	get, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("获取角色资源出错", err)
		return nil
	}
	if get {
		return roleRes
	}

	return nil
}
