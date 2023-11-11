package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

// Group 路由的分组, 例如: v1/user/login
type Group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type Router struct {
	group []*Group
}

// AddRouter 向分组中添加路由
func (g *Group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

// Group 路由中的分组
func (r *Router) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
	}
	r.group = append(r.group, g)

	return g
}

// 2.执行
func (g *Group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	}
}

// Run 1.启动
func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//路径中: account.login (account 组标识,) login 路由标识
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0] // 组标识
		name = strs[1]   //路由标识

	}
	//遍历所有组, 匹配请求中的组标识.
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}
