package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/global"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"sync"
)

var ArmyService = &armyService{}

type armyService struct {
	passBy         sync.RWMutex
	passByPosArmys map[int]map[int]*data.Army // key : posId,armyId
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

func (a *armyService) ScanBlock(roleId int, req *model.ScanBlockReq) ([]model.Army, error) {
	x := req.X
	y := req.Y
	length := req.Length
	out := make([]model.Army, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return out, nil
	}
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	a.passBy.RLock()
	defer a.passBy.RUnlock()

	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			armys, ok := a.passByPosArmys[posId]
			if ok {
				//是否在视野内
				is := armyIsInView(roleId, i, j)
				if is {
					continue
				}
				for _, army := range armys {
					out = append(out, army.ToModel().(model.Army))
				}
			}

		}
	}
	return out, nil
}

func armyIsInView(rid, x, y int) bool {
	return true
}
