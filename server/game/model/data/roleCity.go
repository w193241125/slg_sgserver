package data

import (
	"sync"
	"time"
)

type RoleCity struct {
	mutex      sync.Mutex `xorm:"-"`
	CityId     int        `xorm:"cityId pk autoincr"`
	RId        int        `xorm:"rid"`
	Name       string     `xorm:"name" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	X          int        `xorm:"x"`
	Y          int        `xorm:"y"`
	IsMain     int8       `xorm:"is_main"`
	CurDurable int        `xorm:"cur_durable"`
	CreatedAt  time.Time  `xorm:"created_at"`
	OccupyTime time.Time  `xorm:"occupy_time"`
}

func (m *RoleCity) TableName() string {
	return "map_role_city"
}
