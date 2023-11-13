package main

import (
	"sgserver/config"
	"sgserver/net"
	"sgserver/server/gate"
)

/*
	1. 登陆功能 account.login 通过网关 转发到 登录服务器
	2. 网关如何和登录服务器(websocket)交互 : (是服务器websocket的客户端)
	3. 网关又和游戏客户端交互,(是游戏客户端 websocket的服务端)
	4. websocket 的服务端已经实现了.
	5. websocket 的客户端
	6. 网关: 代理服务器 (代理地址  代理的连接通道)  客户端连接(websocket 连接)
	7. 路由: 接收所有请求(*). 网关的websocket服务端功能
	8. 握手协议 : 检测第一次连接的时候是授信状态(连通, 合法).
*/

func main() {
	host := config.File.MustValue("gate_server", "host", "127.0.0.1")
	port := config.File.MustValue("gate_server", "port", "8004")

	s := net.NewServer(host + ":" + port)
	s.NeedSecret(true)
	gate.Init()
	s.Router(gate.Router)
	s.Start()
	//log.Fatal("登陆服务器成功")
}
