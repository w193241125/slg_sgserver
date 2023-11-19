package logic

import (
	"encoding/json"
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sync"
)

var CoalitionService = &coalitionService{
	unions: make(map[int]*data.Coalition),
}

type coalitionService struct {
	mutex  sync.RWMutex
	unions map[int]*data.Coalition
}

// menbers = [1,2,3,4]
func (c *coalitionService) Load() {
	rr := make([]*data.Coalition, 0)
	err := db.Engine.Table(new(data.Coalition)).Where("state=?", data.UnionRunning).Find(&rr)
	if err != nil {
		log.Println("coalitionService.Load-查询联盟失败 err", err)
	}
	for _, v := range rr {
		c.unions[v.Id] = v
	}
	log.Println("coalitionService.Load-查询联盟成功", len(rr))
}

func (c *coalitionService) List() ([]model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]model.Union, 0)
	for _, v := range c.unions {
		mas := make([]model.Major, 0)
		chairman := v.Chairman
		if role := RoleService.Get(chairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		if role := RoleService.Get(chairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}

		union := v.ToModel().(model.Union)
		union.Major = mas
		uns = append(uns, union)
	}
	return uns, nil
}

func (c *coalitionService) ListCoalition() []*data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]*data.Coalition, 0)
	ms, _ := json.Marshal(c.unions)
	log.Println("c.unions", ms)
	for _, v := range c.unions {
		uns = append(uns, v)
	}
	return uns
}
