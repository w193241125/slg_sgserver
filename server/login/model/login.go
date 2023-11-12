package model

import "time"

const (
	Login = iota
	Logout
)

type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"`
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"`
	Hardware string    `xorm:"hardware"`
}

// xorm 自行指定表名
func (*LoginHistory) TableName() string {
	return "login_history"
}

type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"`
	UId        int       `xorm:"uid"`
	LoginTime  time.Time `xorm:"login_time"`
	LogoutTime time.Time `xorm:"logout_time"`
	Ip         string    `xorm:"ip"`
	Session    string    `xorm:"session"`
	IsLogout   int8      `xorm:"is_logout"`
	Hardware   string    `xorm:"hardware"`
}

// xorm 自行指定表名
func (*LoginLast) TableName() string {
	return "login_last"
}
