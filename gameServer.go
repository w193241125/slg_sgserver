package main

import (
	"sgserver/config"
	"sgserver/net"
	"sgserver/server/game"
)

/*
EnterServerReq
1. 登陆完成, 创角
2. 根据用户查询 此用户所拥有的角色, 没有就 创角
3. 木材, 令牌, 金钱, 主城, 武将 等等...
4. 地图相关的, 城池, 资源土地, 要塞等,需要定义
5. 资源, 军队,城池, 武将等等...需要加载.
*/
type EnterServerReq struct {
	Session string `json:"session"`
}

func main() {
	host := config.File.MustValue("game_server", "host", "127.0.0.1")
	port := config.File.MustValue("game_server", "port", "8001")
	s := net.NewServer(host + ":" + port)
	s.NeedSecret(false)
	game.Init()
	s.Router(game.Router)
	s.Start()
}
