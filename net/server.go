package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	addr   string
	router *Router
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Router(router *Router) {
	s.router = router
}

// Start 启动服务
func (s *Server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http 升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	//允许所有 CORS 跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	//websocket
	//1. http 协议升级为 websocket 协议
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// 打印日志,同时会退出应用程序
		log.Println("websocket 服务连接出错", err)
	}
	log.Println("websocket 服务连接成功")
	//fmt.Println("websocket 服务连接成功")
	//websocket 通道建立之后,不管是客服端还是服务端,都可以收发消息
	//发消息的时候,把消息当做路由来去处理, 消息是有格式的,先定义消息格式
	//客户端 发消息 {Name: "account.login"} 收到之后,进行解析,知道是要处理登录逻辑

	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}
