package data

import (
	"encoding/json"
	"log"
	"sgserver/server/game/model"
	"time"
	"xorm.io/xorm"
)

const (
	UnionDismiss = 0 //解散
	UnionRunning = 1 //运行中
)

type Coalition struct {
	Id           int       `xorm:"id pk autoincr"`
	Name         string    `xorm:"name"`
	Members      string    `xorm:"members"`
	MemberArray  []int     `xorm:"-"`
	CreateId     int       `xorm:"create_id"`
	Chairman     int       `xorm:"chairman"`
	ViceChairman int       `xorm:"vice_chairman"`
	Notice       string    `xorm:"notice"`
	State        int8      `xorm:"state"`
	Ctime        time.Time `xorm:"ctime"`
}

func (c *Coalition) TableName() string {
	return "coalition"
}

// AfterSet xorm提供的方法，在xorm的操作完成后，会调用该方法，将数据转换为我们需要的格式
// 例如：将数据库中的字符串转换为数组
func (c *Coalition) AfterSet(name string, cell xorm.Cell) {
	if name == "members" {
		if cell != nil {
			ss, ok := (*cell).([]uint8) //转为数组
			if ok {
				json.Unmarshal(ss, &c.MemberArray)
			}

			if c.MemberArray == nil {
				c.MemberArray = []int{}
				log.Println("查询联盟后进行数据转换", c.MemberArray)
			}
		}
	}
}

func (c *Coalition) ToModel() interface{} {
	u := model.Union{}
	u.Name = c.Name
	u.Notice = c.Notice
	u.Id = c.Id

	return u
}

type CoalitionApply struct {
	Id      int       `xorm:"id pk autoincr"`
	UnionId int       `xorm:"union_id"`
	RId     int       `xorm:"rid"`
	State   int8      `xorm:"state"`
	Ctime   time.Time `xorm:"ctime"`
}

func (c *CoalitionApply) TableName() string {
	return "coalition_apply"
}

const (
	UnionOpCreate    = 0 //创建
	UnionOpDismiss   = 1 //解散
	UnionOpJoin      = 2 //加入
	UnionOpExit      = 3 //退出
	UnionOpKick      = 4 //踢出
	UnionOpAppoint   = 5 //任命
	UnionOpAbdicate  = 6 //禅让
	UnionOpModNotice = 7 //修改公告
)

type CoalitionLog struct {
	Id       int       `xorm:"id pk autoincr"`
	UnionId  int       `xorm:"union_id"`
	OPRId    int       `xorm:"op_rid"`
	TargetId int       `xorm:"target_id"`
	State    int8      `xorm:"state"`
	Des      string    `xorm:"des"`
	Ctime    time.Time `xorm:"ctime"`
}

func (c *CoalitionLog) TableName() string {
	return "coalition_log"
}
