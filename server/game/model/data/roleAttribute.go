package data

import (
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
