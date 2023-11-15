package logic

import (
	"encoding/json"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/gameConfig/general"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"time"
)

var GeneralService = &generalService{}

type generalService struct {
}

func (g *generalService) GetGenerals(rid int) ([]model.General, error) {
	gs := make([]*data.General, 0)
	gl := &data.General{}
	err := db.Engine.Table(gl).Where("rid=?", rid).Find(&gs)

	if err != nil {
		log.Println("查询武将出错", err)
		return nil, common.New(constant.DBError, "getGenerals武将查询出错")

	}
	if len(gs) <= 0 {
		var count = 0
		for {
			if count >= 3 {
				break
			}
			cfgId := general.General.Rand()
			gen, err := g.NewGeneral(cfgId, rid, 0)
			if err != nil {
				log.Println("生成初始武将失败", err)
				continue
			}

			gs = append(gs, gen)
			count++
		}
	}

	Generals := make([]model.General, 0)
	for _, v := range gs {
		Generals = append(Generals, v.ToModel().(model.General))
	}
	return Generals, nil
}

const (
	GeneralNormal      = 0 //正常
	GeneralComposeStar = 1 //星级合成
	GeneralConvert     = 2 //转换
)

func (g *generalService) NewGeneral(cfgId int, rid int, level int8) (*data.General, error) {
	cfg := general.General.GMap[cfgId]
	//初始武将无技能,但是有技能槽
	sa := make([]*model.GSkill, 3)
	ss, _ := json.Marshal(sa)

	gen := &data.General{
		PhysicalPower: gameConfig.Base.General.PhysicalPowerLimit,
		RId:           rid,
		CfgId:         cfg.CfgId,
		Order:         0,
		CityId:        0,
		Level:         level,
		CreatedAt:     time.Now(),
		CurArms:       cfg.Arms[0],
		HasPrPoint:    0,
		UsePrPoint:    0,
		AttackDis:     0,
		ForceAdded:    0,
		StrategyAdded: 0,
		DefenseAdded:  0,
		SpeedAdded:    0,
		DestroyAdded:  0,
		Star:          cfg.Star,
		StarLv:        0,
		ParentId:      0,
		SkillsArray:   sa,
		Skills:        string(ss),
		State:         GeneralNormal,
	}

	_, err := db.Engine.Table(gen).Insert(gen)

	if err != nil {
		log.Println("初始武将入库失败", err)
		return nil, err
	}
	return gen, nil
}
