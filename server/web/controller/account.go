package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sgserver/constant"
	"sgserver/server/common"
	"sgserver/server/web/logic"
	"sgserver/server/web/model"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {
}

func (a *AccountController) Register(ctx *gin.Context) {
	fmt.Println("--注册")
	/**
	1. 获取请求参数
	2.根据用户名查询数据库是否已有账户, 有 失败, 没有, 注册.
	3.注册成功
	*/
	rq := &model.RegisterReq{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}

	//web服务, 错误格式一般自定义
	err = logic.DefaultAccountLogic.Register(rq)

	if err != nil {
		log.Println("注册业务出错", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil))
}
