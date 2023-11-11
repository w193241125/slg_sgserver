package main

import (
	"fmt"
	"log"
	"sgserver/config"
	"sgserver/net"
	"sgserver/server/login"
)

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	fmt.Println(host)
	port := config.File.MustValue("login_server", "port", "8003")
	fmt.Println(port)

	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	s.Start()
	log.Fatal("登陆服务器成功")
}
