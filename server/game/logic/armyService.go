package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
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

func (r *armyService) GetArmysByCity(rid int, cid int) ([]model.Army, error) {
	mrs := make([]data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=? and cityId=?", rid, cid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]model.Army, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Army))
	}
	return modelMrs, nil
}

func (r *armyService) ScanBlock(req *model.ScanBlockReq) ([]model.Army, error) {
	return nil, nil
}
