package controller

import (
	"sgserver/net"
	"sgserver/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(router *net.Router) {
	g := router.Group("account")
	g.AddRouter("login", a.login)

}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{}
	loginRes.UId = 1
	loginRes.Username = "admin"
	loginRes.Session = "as"
	loginRes.Password = ""
	rsp.Body.Msg = loginRes
}
