package data

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
)

type RoleRes struct {
	Id     int `xorm:"id pk autoincr"`
	RId    int `xorm:"rid"`
	Wood   int `xorm:"wood"`
	Iron   int `xorm:"iron"`
	Stone  int `xorm:"stone"`
	Grain  int `xorm:"grain"`
	Gold   int `xorm:"gold"`
	Decree int `xorm:"decree"` //令牌
}

func (r *RoleRes) TableName() string {
	return "role_res"
}

func (r *RoleRes) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	yield := GetYield(r.RId)
	p.GoldYield = yield.Gold
	p.GrainYield = yield.Grain
	p.StoneYield = yield.Stone
	p.IronYield = yield.Iron
	p.WoodYield = yield.Wood
	p.DepotCapacity = 10000
	return p
}

// 资源产量
type Yield struct {
	Wood  int
	Iron  int
	Stone int
	Grain int
	Gold  int
}

var RoleResDao = &roleResDao{
	rrChan: make(chan *RoleRes, 100),
}

func init() {
	go RoleResDao.run()
}

func (r *roleResDao) run() {
	for {
		select {
		case rr := <-r.rrChan:
			//更新操作
			_, err := db.Engine.Table(new(RoleRes)).ID(rr.Id).Cols("wood", "iron", "stone", "grain", "gold").Update(rr)
			if err != nil {
				log.Println("更新角色资源失败", err)
			}

		}
	}
}

type roleResDao struct {
	rrChan chan *RoleRes
}

func (r *RoleRes) SyncExecute() {
	RoleResDao.rrChan <- r
}
