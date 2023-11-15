package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
)

type warservice struct {
}

var WarService = &warservice{}

func (w *warservice) GetWarReport(rid int) ([]model.WarReport, error) {
	wr := make([]data.WarReport, 0)
	wdb := &data.WarReport{}
	err := db.Engine.Table(wdb).Where("a_rid=? or d_rid=?", rid, rid).Limit(30, 0).Desc("ctime").Find(&wr)
	if err != nil {
		log.Println("查询战报失败")
		return nil, common.New(constant.DBError, "查询战报出错")
	}
	modelWr := make([]model.WarReport, 0)
	for _, v := range wr {
		modelWr = append(modelWr, v.ToModel().(model.WarReport))
	}
	return modelWr, nil
}
