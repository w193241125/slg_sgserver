package logic

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var RoleBuildService = &roleBuildService{}

type roleBuildService struct {
}

func (r *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	builds := make([]data.MapRoleBuild, 0)
	build := &data.MapRoleBuild{}
	err := db.Engine.Table(build).Where("rid=?", rid).Find(&builds)

	modelBuilds := make([]model.MapRoleBuild, 0)
	if err != nil {
		log.Println("查询建筑出错", err)
		return modelBuilds, err
	}
	for _, v := range builds {
		modelBuilds = append(modelBuilds, v.ToModel().(model.MapRoleBuild))
	}
	return modelBuilds, nil
}
