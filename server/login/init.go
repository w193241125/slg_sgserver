package login

import (
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/login/controller"
)

var Router = &net.Router{}

func Init() {
	//测试数据库,并且初始化
	db.TestDB()

	//还有别的init 方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
