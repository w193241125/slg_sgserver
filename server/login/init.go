package login

import (
	"sgserver/net"
	"sgserver/server/login/controller"
)

var Router = &net.Router{}

func Init() {
	//还有别的init 方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
