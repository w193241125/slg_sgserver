package data

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/model"
	"time"
)

type RoleAttribute struct {
	Id              int            `xorm:"id pk autoincr"`
	RId             int            `xorm:"rid"`
	UnionId         int            `xorm:"-"`                 //联盟id
	ParentId        int            `xorm:"parent_id"`         //上级id（被沦陷）
	CollectTimes    int8           `xorm:"collect_times"`     //征收次数
	LastCollectTime time.Time      `xorm:"last_collect_time"` //最后征收的时间
	PosTags         string         `xorm:"pos_tags"`          //位置标记
	PosTagArray     []model.PosTag `xorm:"-"`                 //上面的位置标记都是存在一条字符串里面, 使用的时候就会转成一个数组来使用. 所以单独拎出来, 做一个模型.代表对应的一块位置(坐标点,名称).
}

func (r *RoleAttribute) TableName() string {
	return "role_attribute"
}

var RoleAttrDao = &roleAttrDao{
	raChan: make(chan *RoleAttribute, 100),
}

type roleAttrDao struct {
	raChan chan *RoleAttribute
}

func init() {
	go RoleAttrDao.run()
}

func (r *roleAttrDao) run() {
	for {
		select {
		case ra := <-r.raChan:
			upd, err := db.Engine.Table(new(RoleAttribute)).ID(ra.Id).Cols("parent_id", "collect_times", "last_collect_time", "pos_tags").Update(ra)
			if err != nil || upd == 0 {
				log.Println("角色征收信息更新失败", err, upd)
			}
		}
	}
}
func (r *RoleAttribute) SyncExecute() {
	RoleAttrDao.raChan <- r
}
