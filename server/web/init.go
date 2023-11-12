package web

import (
	"github.com/gin-gonic/gin"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/web/controller"
	"sgserver/server/web/middleware"
)

var Router = &net.Router{}

func Init(router *gin.Engine) {
	//测试数据库,并且初始化
	db.TestDB()

	//还有别的init 方法
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	router.Use(middleware.Cors()) // 跨域中间件
	router.Any("/account/register", controller.DefaultAccountController.Register)
}
