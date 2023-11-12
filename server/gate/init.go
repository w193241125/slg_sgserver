package gate

import (
	"sgserver/net"
	"sgserver/server/gate/controller"
)

var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	controller.GateHander.Router(Router)
}
