package logic

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var ArmyService = &armyService{}

type armyService struct {
}

func (r *armyService) GetArmys(rid int) ([]model.Army, error) {
	builds := make([]data.Army, 0)
	build := &data.Army{}
	err := db.Engine.Table(build).Where("rid=?", rid).Find(&builds)

	armys := make([]model.Army, 0)
	if err != nil {
		log.Println("查询军队出错", err)
		return armys, err
	}
	for _, v := range builds {
		armys = append(armys, v.ToModel().(model.Army))
	}
	return armys, nil
}
