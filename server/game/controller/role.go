package controller

import "sgserver/net"

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//进入游戏逻辑
	// Session 是否合法, 合法可取出登录用户的ID
	//根据用户ID查询对应的游戏角色..
	//根据角色ID 查询角色拥有的资源, 有就返回, 没有就初始化.

}
