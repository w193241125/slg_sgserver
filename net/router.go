package net

import (
	"log"
	"strings"
	"sync"
)

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)
type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

// Group 路由的分组, 例如: v1/user/login
type Group struct {
	mutex         sync.RWMutex
	prefix        string
	handlerMap    map[string]HandlerFunc
	middlewareMap map[string][]MiddlewareFunc //针对某个路由的中间件
	middlewares   []MiddlewareFunc            // 针对组的中间件
}

type Router struct {
	group []*Group
}

// AddRouter 向分组中添加路由
// Group类型中的AddRouter方法用于向Group实例添加路由。
//
// 参数：
//   - name: 路由名字符串，用于唯一标识该路由。
//   - handlerFunc: 处理函数，用于处理该路由的请求。
//   - middlewares: 可变参数，中间件函数列表，用于在处理路由请求前或后执行。
//
// 返回值：无
func (g *Group) AddRouter(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.handlerMap[name] = handlerFunc
	g.middlewareMap[name] = middlewares
}

// Use 使用中间件
func (g *Group) Use(middlewares ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Group 路由中的分组
func (r *Router) Group(prefix string) *Group {
	g := &Group{
		prefix:        prefix,
		handlerMap:    make(map[string]HandlerFunc),
		middlewareMap: make(map[string][]MiddlewareFunc),
	}
	r.group = append(r.group, g)

	return g
}

// 2.执行
func (g *Group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h, ok := g.handlerMap[name]
	if !ok {
		h, ok = g.handlerMap["*"]
		if !ok {
			log.Println("路由未定义!")
		}
	}
	if ok {
		for i := 0; i < len(g.middlewares); i++ {
			h = g.middlewares[i](h)
		}
		mm, ok := g.middlewareMap[name]
		if ok {
			for i := 0; i < len(mm); i++ {
				h = mm[i](h)
			}
		}
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
		} else if g.prefix == "*" {
			g.exec(name, req, rsp)
		}
	}
}
