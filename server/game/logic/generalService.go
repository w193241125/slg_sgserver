package logic

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

var GeneralService = &generalService{}

type generalService struct {
}

func (r *generalService) GetGenerals(rid int) ([]model.General, error) {
	builds := make([]data.General, 0)
	build := &data.General{}
	err := db.Engine.Table(build).Where("rid=?", rid).Find(&builds)

	Generals := make([]model.General, len(builds))
	if err != nil {
		log.Println("查询武将出错", err)
		return Generals, err
	}
	for _, v := range builds {
		Generals = append(Generals, v.ToModel().(model.General))
	}
	return Generals, nil
}
